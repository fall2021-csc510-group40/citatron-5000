package main

import (
	"core/schema"
	"core/source"
)

func main() {
	var searches []source.Search
	searches = append(searches, source.SourceSearchACM)

	for _, search := range searches {
		search(&schema.Work{
			Title: "test",
		})
	}
}
