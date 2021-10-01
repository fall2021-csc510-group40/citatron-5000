/*MIT License

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

// Package server implements the REST API used by front-ends to query the works and format them
package server

import (
	"context"
	"core/db"
	"core/formatter"
	"core/schema"
	"encoding/json"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// MaxResults is the maximum results for a single query
const MaxResults = 100

type Config struct {
	Database *db.DatabaseConfig `json:"database"`
}

// Server represents a generic citation server
type Server struct {
	DB *db.Database
}

// New creates a new citation server
func New(config *Config) (*Server, error) {
	database, err := db.New(config.Database)
	if err != nil {
		return nil, err
	}
	return &Server{
		database,
	}, nil
}

// ListenAndServe starts the citation server on the given address
func (s *Server) ListenAndServe(addr string) error {
	http.HandleFunc("/search", s.search)
	http.HandleFunc("/format", s.format)
	return http.ListenAndServe(addr, nil)
}

func contextFromRequest(req *http.Request) context.Context {
	ctx := context.Background()
	if req_id := req.Header.Get("X-Request-ID"); req_id != "" {
		ctx = context.WithValue(ctx, "req_id", req_id)
	}

	return ctx
}

// search handles an HTTP request to search for works
func (s *Server) search(w http.ResponseWriter, req *http.Request) {
	ctx := contextFromRequest(req)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var searchRequest schema.SearchRequest
	if err := json.Unmarshal(body, &searchRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var searchResponse schema.SearchResponse

	results, err := s.DB.Search(ctx, searchRequest.Query)
	if err == nil {
		if len(results) > MaxResults {
			results = results[0:MaxResults]
		}

		searchResponse.Results = results
	} else {
		searchResponse.Error = err.Error()
		log.Errorf("search failed: %v", err)
	}

	resp, err := json.Marshal(searchResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// format handles an HTTP request to format a work into a citation
func (s *Server) format(w http.ResponseWriter, req *http.Request) {
	ctx := contextFromRequest(req)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error while reading format request body: %v", err)
		return
	}

	var formatRequest schema.FormatRequest
	if err := json.Unmarshal(body, &formatRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error while parsing format request body: %v", err)
		return
	}

	var formatResponse schema.FormatResponse

	work := formatRequest.Work
	if work.ID != "" {
		work, err = s.DB.GetWorkById(ctx, work.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Cannot get work by id: %v", err)
			return
		}
	}

	if formatResponse.Error == "" {
		switch formatRequest.Format {
		case "bibtex":
			formatResponse.Result = formatter.BibtexFormat(work)
		default:
			formatResponse.Result = formatter.PlaintextFormat(work)
		}
	}

	resp, err := json.Marshal(formatResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
