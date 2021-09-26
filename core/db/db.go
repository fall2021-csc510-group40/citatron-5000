package db

import (
	"core/schema"
	"core/source"
	"core/util"
	"sort"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
)

// SimilarityThreshold is the max string distance for a search result from the target
const SimliartyThreshold = 10

var searches = [...]source.Search{source.SourceSearchACM, source.SourceSearchCrossRef}

// Database represents a generic database instance for works
type Database struct {
}

// New constructs a new database
func New() *Database {
	return &Database{}
}

// Search searches the database or search sources for a given work
func (d *Database) Search(work *schema.Work) ([]*schema.Work, error) {
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

	// Limit baesd on similarity threshold
	i := 0
	for ; i < len(works); i++ {
		if workToDistance[works[i]] > SimliartyThreshold {
			break
		}
	}

	return works[0:i], nil
}
