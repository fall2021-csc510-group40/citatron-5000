package server

import (
	"core/db"
	"core/schema"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// MaxResults is the maximum results for a single query
const MaxResults = 10

// Server represents a generic citation server
type Server struct {
	DB *db.Database
}

// New creates a new citation server
func New() *Server {
	server := &Server{}
	server.DB = db.New()
	return server
}

// ListenAndServe starts the citation server on the given address
func (s *Server) ListenAndServe(addr string) error {
	http.HandleFunc("/search", s.search)
	http.HandleFunc("/format", s.format)
	return http.ListenAndServe(addr, nil)
}

// search handles an HTTP request to search for works
func (s *Server) search(w http.ResponseWriter, req *http.Request) {
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

	results, err := s.DB.Search(searchRequest.Query)
	if err == nil {
		if len(results) > MaxResults {
			results = results[0:MaxResults]
		}

		searchResponse.Results = results
	} else {
		searchResponse.Error = err.Error()
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

}
