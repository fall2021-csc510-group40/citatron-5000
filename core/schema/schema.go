package schema

import (
	"core/util"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Work struct {
	ID   string
	Hash string

	Type string

	DOI   string
	Arxiv string
	ISBN  string

	Title   string
	Authors []string

	Version string
	Venue   string
	Page    string

	Year  int
	Month int
	Day   int

	Keywords []string
}

type SearchRequest struct {
	Query *Work
}

type SearchResponse struct {
	Results []*Work
	Error   string
}

type FormatRequest struct {
	ID     string
	Work   *Work
	Format string
}

type FormatResponse struct {
	Result string
	Error  string
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
