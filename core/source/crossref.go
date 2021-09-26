package source

import (
	"core/schema"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SourceSearchCrossRef searches CrossRef for works
func SourceSearchCrossRef(titleGetter *schema.Work) ([]*schema.Work, error) {
	var works []*schema.Work
	var title string = titleGetter.Title
	query := url.Values{}
	query.Add("query.bibliographic", title)
	query.Add("rows", "5")

	resp, err := http.Get("http://api.crossref.org/works?" + query.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Status  string `json:"status"`
		Message struct {
			Items []struct {
				DOI    string   `json:"DOI"`
				Title  []string `json:"title"`
				Author []struct {
					Given  string `json:"given"`
					Family string `json:"family"`
				} `json:"author"`
				EditionNumber string `json:"edition_number"`
				Publisher     string `json:"publisher"`
				Created       struct {
					DateParts [][]int `json:"date-parts"`
				} `json:"created"`
				Type string `json:"type"`
				Page string `json:"page"`
			} `json:"items"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	for _, i := range parsed.Message.Items {
		w := &schema.Work{}
		w.DOI = i.DOI
		if len(i.Title) > 0 {
			w.Title = i.Title[0]
		} else {
			w.Title = ""
		}

		for _, j := range i.Author {
			w.Authors = append(w.Authors, (j.Given + " " + j.Family))
		}
		w.Version = i.EditionNumber
		w.Venue = i.Publisher
		if len(i.Created.DateParts) > 0 {
			w.Year = i.Created.DateParts[0][0]
			w.Month = i.Created.DateParts[0][1]
			w.Day = i.Created.DateParts[0][2]
		} else {
			w.Year = 0
			w.Month = 0
			w.Day = 0
		}

		w.Type = i.Type
		w.Page = i.Page
		works = append(works, w)
	}

	return works, nil

}
