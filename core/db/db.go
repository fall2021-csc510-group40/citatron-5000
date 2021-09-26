package db

import (
	"context"
	"core/schema"
	"core/source"
	"core/util"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimilarityThreshold is the max string distance for a search result from the target
const SimliartyThreshold = 10

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

func (d *Database) SearchDatabase(work *schema.Work) ([]*schema.Work, error) {
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
func (d *Database) Search(work *schema.Work) ([]*schema.Work, error) {
	return d.SearchDatabase(work)
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

	// Sort works based on title string distance
	var works []*schema.Work
	workToDistance := map[*schema.Work]int{}

	for _, w := range uniqueWorks {
		works = append(works, w)
		workToDistance[w] = levenshtein.ComputeDistance(work.Title, strings.ToLower(w.Title))
	}

	sort.Slice(works, func(i, j int) bool {
		return workToDistance[works[i]] < workToDistance[works[j]]
	})

	// Limit based on similarity threshold
	i := 0
	for ; i < len(works); i++ {
		if workToDistance[works[i]] > SimliartyThreshold {
			break
		}
	}

	return works[0:i], nil
}

