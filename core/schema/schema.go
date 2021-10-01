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
package schema

import (
	"core/util"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

// Work represents a generic work like article, book, etc.
type Work struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Hash string `json:"hash" bson:",omitempty"`

	Type string `json:"type" bson:",omitempty"`

	DOI   string `json:"doi" bson:",omitempty"`
	Arxiv string `json:"arxiv" bson:",omitempty"`
	ISBN  string `json:"isbn" bson:",omitempty"`

	Title   string   `json:"title" bson:",omitempty"`
	Authors []string `json:"authors" bson:",omitempty"`

	Version string `json:"version" bson:",omitempty"`
	Venue   string `json:"venue" bson:",omitempty"`
	Page    string `json:"page" bson:",omitempty"`

	Year  int `json:"year" bson:",omitempty"`
	Month int `json:"month" bson:",omitempty"`
	Day   int `json:"day" bson:",omitempty"`

	Keywords []string `json:"keywords" bson:",omitempty"`
}

// SearchRequest is the body of a search request
type SearchRequest struct {
	Query *Work `json:"query"`
}

// SearchResponse is the body of a search response
type SearchResponse struct {
	Results []*Work `json:"results"`
	Error   string  `json:"error"`
}

// FormatRequest is the body of a format request
type FormatRequest struct {
	Work   *Work  `json:"work"`
	Format string `json:"format"`
}

// FormatResponse is the body of a format response
type FormatResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// Normalize normalizes the data for a work and then populates its hash
func (w *Work) Normalize() error {
	// Clean strings
	w.Type = util.CleanString(w.Type)

	w.DOI = util.CleanString(w.DOI)
	w.Arxiv = util.CleanString(w.Arxiv)
	w.ISBN = util.CleanString(w.ISBN)
	w.Title = util.CleanString(w.Title)

	for i, v := range w.Authors {
		w.Authors[i] = util.CleanString(v)
	}

	w.Version = util.CleanString(w.Version)
	w.Venue = util.CleanString(w.Venue)
	w.Page = util.CleanString(w.Page)

	for i, v := range w.Keywords {
		w.Keywords[i] = util.CleanString(v)
	}

	// Required fields
	if w.Title == "" {
		return errors.New("no title")
	}

	if len(w.Authors) == 0 {
		return errors.New("no author")
	}

	// Alphabetize keywords
	sort.Strings(w.Keywords)

	// Calculate hash
	h := sha256.New()

	var data []string
	data = append(data, util.RemoveAllPunctuation(strings.ToLower(w.Authors[0])))
	data = append(data, util.RemoveAllPunctuation(strings.ToLower(w.Title)))

	for _, d := range data {
		h.Write([]byte(d))
	}

	w.Hash = hex.EncodeToString(h.Sum(nil))
	return nil
}

// Coalesce merges data with this work and another
func (w *Work) Coalesce(other *Work) {
	// Coalesce fields
	if w.Type == "" {
		w.Type = other.Type
	}

	if w.DOI == "" {
		w.DOI = other.DOI
	}

	if w.Arxiv == "" {
		w.Arxiv = other.Arxiv
	}

	if w.ISBN == "" {
		w.ISBN = other.ISBN
	}

	if w.Version == "" {
		w.Version = other.Version
	}

	if w.Venue == "" {
		w.Venue = other.Venue
	}

	if w.Page == "" {
		w.Page = other.Page
	}

	// TODO: consider merging the whole date
	if w.Day == 0 {
		w.Day = other.Day
	}

	if w.Month == 0 {
		w.Month = other.Month
	}

	if w.Year == 0 {
		w.Year = other.Year
	}
}
