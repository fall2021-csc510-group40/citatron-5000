/*
MIT License

Copyright (c) 2021 fall2021-csc510-group40

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package db

import (
	"context"
	"core/schema"
	"core/source"
	"core/util"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
)

var searches = [...]source.Search{source.SearchACM, source.SearchCrossRef}

type DatabaseConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`

	Username string `json:"username"`
	Password string `json:"password"`

	SourceTimeout   util.Duration `json:"source_timeout"`
	DatabaseTimeout util.Duration `json:"database_timeout"`

	LogLevel string `json:"log_level"`
}

// Database represents a generic database instance for works
type Database struct {
	*mongo.Database

	config *DatabaseConfig

	logger log.FieldLogger
}

// New constructs a new database
func New(config *DatabaseConfig) (*Database, error) {
	credential := options.Credential{
		Username: config.Username,
		Password: config.Password,
		AuthSource: config.Database,
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.URI).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New()
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(level)

	return &Database{
		Database: client.Database(config.Database),
		config:   config,
		logger:   logger.WithField("component", "db"),
	}, nil
}

// ExactSearch performs a database search for all works that exactly match the fields provided in the work argument
func (d *Database) ExactSearch(ctx context.Context, work *schema.Work) (res []*schema.Work, err error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.DatabaseTimeout.Duration)
	defer cancel()

	coll := d.Collection("works")
	query, err := convertToDoc(work)
	if err != nil {
		return
	}

	cur, err := coll.Find(ctx, query)
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	if err != nil {
		return
	}

	return
}

// GetWorkById returns the work from the database with the idString as the id.
func (d *Database) GetWorkById(ctx context.Context, idString string) (w *schema.Work, err error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.DatabaseTimeout.Duration)
	defer cancel()

	coll := d.Collection("works")
	id, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return
	}

	query := bson.M{
		"_id": id,
	}

	res := coll.FindOne(ctx, query)
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

// TextSearch performs text search on the work title using mongodb's build-in text index.
// The results are ordered by score provided by the database
func (d *Database) TextSearch(ctx context.Context, work *schema.Work) ([]*schema.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.DatabaseTimeout.Duration)
	defer cancel()

	collection := d.Collection("works")
	query := bson.M{
		"$text": bson.M{
			"$search": work.Title,
		},
	}

	var results []*schema.Work

	cur, err := collection.Find(ctx, query, &options.FindOptions{
		Sort: bson.M{
			"score": bson.M{
				"$meta": "textScore",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// SearchSources searches the database or search sources for a given work
func (d *Database) SearchSources(ctx context.Context, work *schema.Work) ([]*schema.Work, error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.SourceTimeout.Duration)
	defer cancel()

	logger := d.logger
	if req_id := ctx.Value("req_id"); req_id != nil {
		logger = logger.WithField("req_id", req_id)
	}

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

			works, err := s(ctx, work)
			if err != nil {
				logger.Errorf("Error from source: %v", err)
				return
			}
			logger.Debugf("Source responded")


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

func (d *Database) getWorksByHashes(ctx context.Context, h []string) (works []*schema.Work, err error) {
	ctx, cancel := context.WithTimeout(ctx, d.config.DatabaseTimeout.Duration)
	defer cancel()

	coll := d.Collection("works")
	query := bson.M{
		"hash": bson.M{
			"$in": h,
		},
	}

	cur, err := coll.Find(ctx, query)
	if err != nil {
		return
	}

	err = cur.All(ctx, &works)
	return
}

func (d *Database) rewriteMultipleWorks(ctx context.Context, works []*schema.Work) error {
	ctx, cancel := context.WithTimeout(ctx, d.config.DatabaseTimeout.Duration)
	defer cancel()

	logger := d.logger
	if req_id := ctx.Value("req_id"); req_id != nil {
		logger = logger.WithField("req_id", req_id)
	}

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

	logger.Debugf("Bulk write models: %v", models)

	coll := d.Collection("works")
	_, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	return nil
}

// Search looks for the documents in external sources and updates the database accordingly.
// If an exact match exists in the database, all the works matching all the fields in the request are returned.
// Otherwise, the external search is performed first, and the results are uploaded to the database. Then the text search
// is performed using the database's build-in index and the results are returned to the user.
func (d *Database) Search(ctx context.Context, work *schema.Work) (works []*schema.Work, err error) {
	logger := d.logger
	if req_id := ctx.Value("req_id"); req_id != nil {
		logger = logger.WithField("req_id", req_id)
	}

	works, err = d.ExactSearch(ctx, work)
	if err != nil || len(works) == 0 {
		logger.Debugf("Exact search failed: %v", err)
	} else {
		return
	}

	candidates, err := d.SearchSources(ctx, work)
	if err != nil || len(candidates) == 0 {
		logger.Debugf("Sources search failed: %v", err)
	} else {
		logger.Debugf("Found %v candidates from the sources", len(candidates))
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

		dbCandidates, err := d.getWorksByHashes(ctx, hashes)
		if err != nil {
			return nil, err
		}

		logger.Debugf("Got %v matches in the database", len(dbCandidates))

		for _, dbCandidate := range dbCandidates {
			sourceCandidate := candidateByHash[dbCandidate.Hash]
			dbCandidate.Coalesce(sourceCandidate)
			candidateByHash[dbCandidate.Hash] = dbCandidate
		}

		var worksToWrite []*schema.Work
		for _, w := range candidateByHash {
			worksToWrite = append(worksToWrite, w)
		}

		logger.Debugf("Writing candidates to database")
		err = d.rewriteMultipleWorks(ctx, worksToWrite)
		if err != nil {
			logger.Errorf("Unable to write candidates to the database: %v", err)
			return nil, err
		}
	}

	return d.TextSearch(ctx, work)
}
