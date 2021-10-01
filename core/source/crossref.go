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
package source

import (
	"context"
	"core/schema"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SearchCrossRef searches CrossRef for works
func SearchCrossRef(ctx context.Context, w *schema.Work) ([]*schema.Work, error) {
	var works []*schema.Work
	title := w.Title
	query := url.Values{}
	query.Add("query.bibliographic", title)
	query.Add("rows", "5")

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://api.crossref.org/works?" + query.Encode(), nil)
	resp, err := http.DefaultClient.Do(req)
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

		if len(i.Created.DateParts) > 0 && len(i.Created.DateParts[0]) > 2 {
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
