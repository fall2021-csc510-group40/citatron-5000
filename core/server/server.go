package server

import (
	"core/db"
	"core/schema"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Server struct {
	DB *db.Database
}

func New() *Server {
	server := &Server{}
	server.DB = db.New()
	return server
}

func (s *Server) ListenAndServe(addr string) error {
	http.HandleFunc("/search", s.search)
	http.HandleFunc("/format", s.format)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) search(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (s *Server) format(w http.ResponseWriter, req *http.Request) {

}
