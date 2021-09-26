package db

import (
	"core/schema"
	"core/source"
)

var searches = [...]source.Search{source.SourceSearchACM, source.SourceSearchCrossRef}

type Database struct {
}

func New() *Database {
	return &Database{}
}

func (d *Database) Search(work *schema.Work) ([]*schema.Work, error) {
	uniqueWorks := map[string]*schema.Work{}

	for _, s := range searches {
		works, err := s(work)
		if err != nil {
			return nil, err
		}

		for _, w := range works {
			if err := w.Normalize(); err != nil {
				continue
			}

			if v, ok := uniqueWorks[w.Hash]; ok {
				v.Coalesce(w)
			} else {
				uniqueWorks[w.Hash] = w
			}
		}
	}

	var works []*schema.Work
	for _, w := range uniqueWorks {
		works = append(works, w)
	}

	return works, nil
}
