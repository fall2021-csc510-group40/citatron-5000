package schema

import (
	"core/util"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Work struct {
	ID   string `json:"id"`
	Hash string `json:"hash"`

	Type string `json:"type"`

	DOI   string `json:"doi"`
	Arxiv string `json:"arxiv"`
	ISBN  string `json:"isbn"`

	Title   string   `json:"title"`
	Authors []string `json:"authors"`

	Version string `json:"version"`
	Venue   string `json:"venue"`
	Page    string `json:"page"`

	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`

	Keywords []string `json:"keywords"`
}

type SearchRequest struct {
	Query *Work `json:"query"`
}

type SearchResponse struct {
	Results []*Work `json:"results"`
	Error   string  `json:"error"`
}

type FormatRequest struct {
	ID     string `json:"id"`
	Work   *Work  `json:"work"`
	Format string `json:"format"`
}

type FormatResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

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

	if w.Year == 0 {
		return errors.New("no year")
	}

	// Alphabetize authors and keywords
	sort.Strings(w.Authors)
	sort.Strings(w.Keywords)

	// Calculate hash
	h := sha256.New()

	var data []string
	data = append(data, fmt.Sprintf("%d", w.Year))
	data = append(data, util.RemoveAllPunctuation(strings.ToLower(w.Title)))

	for _, d := range data {
		h.Write([]byte(d))
	}

	w.Hash = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (w *Work) Coalesce(other *Work) {
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

	if len(w.Authors) == 0 {
		w.Authors = other.Authors
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

	if w.Day == 0 {
		w.Day = other.Day
	}

	if w.Month == 0 {
		w.Month = other.Month
	}
}
