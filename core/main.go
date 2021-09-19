package main

import (
	"core/schema"
	"core/source"
	"fmt"
)

func main() {
	var searches []source.Search
	searches = append(searches, source.SourceSearchACM)

	for _, search := range searches {
		results, err := search(&schema.Work{Title: "test"})
		if err != nil {
			panic(err)
		}

		for _, w := range results {
			fmt.Println(w.Type)
			fmt.Println(w.Title)
			fmt.Println(w.Venue)
			fmt.Println(w.Page)
			fmt.Println(w.DOI)
			fmt.Println(w.Day, w.Month, w.Year)

			for _, a := range w.Authors {
				fmt.Println("\t" + a)
			}
		}
	}
}
