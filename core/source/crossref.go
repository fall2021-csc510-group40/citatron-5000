package source

import (
	"core/schema"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Response struct {
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
			Event         struct {
				Start struct {
					DateParts [][]int `json:"date-parts"`
				} `json:"start"`
			} `json:"event"`
			Type string `json:"type"`
			Page string `json:"page"`
		} `json:"items"`
	} `json:"message"`
}

func SourceSearchCrossRef(titleGetter *schema.Work) ([]*schema.Work, error) {
	var works []*schema.Work
	var title string = titleGetter.Title
	query := url.Values{}
	query.Add("query.bibliographic", title)
	query.Add("rows", "5")

	resp, err := http.Get("http://api.crossref.org/works?" + query.Encode())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	var parsed Response

	if err := json.Unmarshal(body, &parsed); err != nil {
		panic(err)
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
		if len(i.Event.Start.DateParts) > 0 {
			w.Year = i.Event.Start.DateParts[0][0]
			w.Month = i.Event.Start.DateParts[0][1]
			w.Day = i.Event.Start.DateParts[0][2]
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
