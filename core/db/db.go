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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func convertToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

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
		client.Database("works"),
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

func (d *Database) TextSearch(work *schema.Work) ([]*schema.Work, error) {
	collection := d.Collection("works")
	query := bson.M{
		"$text": bson.M{
			"$search": work.Title,
		},
	}

	var results []*schema.Work

	log.Printf("Searching for %s", work.Title)
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

// Search searches the database or search sources for a given work
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
	return uniqueWorks, nil
}
