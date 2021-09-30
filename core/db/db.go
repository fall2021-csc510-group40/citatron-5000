/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package db

import (
	"context"
	"core/schema"
	"core/source"
	"core/util"
	"log"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var searches = [...]source.Search{source.SourceSearchACM, source.SourceSearchCrossRef}

// Database represents a generic database instance for works
type Database struct {
	*mongo.Database
}

// New constructs a new database
func New(uri string) (*Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Database{
		client.Database("citeman"),
	}, nil
}

func (d *Database) ExactSearch(work *schema.Work) (res []*schema.Work, err error) {
	coll := d.Collection("works")
	query, err := convertToDoc(work)
	if err != nil {
		return
	}

	cur, err := coll.Find(context.Background(), query)
	if err != nil {
		return
	}

	err = cur.All(context.Background(), &res)
	if err != nil {
		return
	}

	return
}

func (d *Database) GetWorkById(idString string) (w *schema.Work, err error) {
	coll := d.Collection("works")
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return
	}

	query := bson.M{
		"_id": id,
	}

	res := coll.FindOne(context.Background(), query)
	w = &schema.Work{}
	err = res.Decode(w)
	return
}

func convertToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func (d *Database) TextSearch(work *schema.Work) ([]*schema.Work, error) {
	collection := d.Collection("works")
	query := bson.M{
		"$text": bson.M{
			"$search": work.Title,
		},
	}

	var results []*schema.Work

	cur, err := collection.Find(context.Background(), query, &options.FindOptions{
		Sort: bson.M{
			"score": bson.M{
				"$meta": "textScore",
			},
		},
	})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = cur.All(context.Background(), &results)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return results, nil
}

// SearchSources searches the database or search sources for a given work
func (d *Database) SearchSources(work *schema.Work) ([]*schema.Work, error) {
	// Normalize search data
	work.Title = strings.ToLower(util.CleanString(work.Title))

	// Find all works
	var wg sync.WaitGroup
	var mu sync.Mutex
	uniqueWorks := map[string]*schema.Work{}

	for _, s := range searches {
		wg.Add(1)

		go func(s source.Search) {
			defer wg.Done()

			works, err := s(work)
			if err != nil {
				log.Printf("Error from source: %v", err)
				return
			}

			for _, w := range works {
				if err := w.Normalize(); err != nil {
					continue
				}

				mu.Lock()

				if v, ok := uniqueWorks[w.Hash]; ok {
					v.Coalesce(w)
				} else {
					uniqueWorks[w.Hash] = w
				}

				mu.Unlock()
			}
		}(s)
	}

	wg.Wait()

	var works []*schema.Work
	for _, w := range uniqueWorks {
		works = append(works, w)
	}

	return works, nil
}

func (d *Database) getWorksByHashes(h []string) (works []*schema.Work, err error) {
	coll := d.Collection("works")
	query := bson.M{
		"hash": bson.M{
			"$in": h,
		},
	}

	cur, err := coll.Find(context.Background(), query)
	if err != nil {
		return
	}

	err = cur.All(context.Background(), &works)
	return
}

func (d *Database) rewriteMultipleWorks(works []*schema.Work) error {
	var models []mongo.WriteModel

	for _, w := range works {
		doc, err := convertToDoc(w)
		if err != nil {
			return err
		}

		m := doc.Map()
		delete(m, "_id")
		update := bson.M{
			"$set": m,
		}

		model := mongo.
			NewUpdateOneModel().
			SetUpdate(update).
			SetFilter(bson.M{
				"hash": w.Hash,
			}).
			SetUpsert(true)

		models = append(models, model)
	}

	log.Printf("Bulk write models: %v", models)

	coll := d.Collection("works")
	_, err := coll.BulkWrite(context.Background(), models)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Search(work *schema.Work) (works []*schema.Work, err error) {
	works, err = d.ExactSearch(work)
	if err != nil || len(works) == 0 {
		log.Printf("Exact search failed: %v", err)
	} else {
		return
	}

	candidates, err := d.SearchSources(work)
	if err != nil || len(candidates) == 0 {
		log.Printf("Sources search failed: %v", err)
	} else {
		log.Printf("Found %v candidates from the sources", len(candidates))
		var hashes []string
		candidateByHash := make(map[string]*schema.Work)
		for _, c := range candidates {
			err = c.Normalize()
			if err != nil {
				return
			}
			if other, ok := candidateByHash[c.Hash]; ok {
				other.Coalesce(c)
			} else {
				candidateByHash[c.Hash] = c
			}
			hashes = append(hashes, c.Hash)
		}

		dbCandidates, err := d.getWorksByHashes(hashes)
		if err != nil {
			return nil, err
		}

		log.Printf("Got %v matches in the database", len(dbCandidates))

		for _, dbCandidate := range dbCandidates {
			sourceCandidate := candidateByHash[dbCandidate.Hash]
			dbCandidate.Coalesce(sourceCandidate)
			candidateByHash[dbCandidate.Hash] = dbCandidate
		}

		var worksToWrite []*schema.Work
		for _, w := range candidateByHash {
			worksToWrite = append(worksToWrite, w)
		}

		log.Printf("Writing candidates to database")
		err = d.rewriteMultipleWorks(worksToWrite)
		if err != nil {
			log.Printf("Unable to write candidates to the database: %v", err)
			return nil, err
		}
	}

	return d.TextSearch(work)
}
