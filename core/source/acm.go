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
	"core/util"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getAcmWorks parses the raw HTML from an ACM result page
func getAcmWorks(doc *goquery.Document) []*schema.Work {
	var works []*schema.Work

	doc.Find(".issue-item").Each(func(i int, s *goquery.Selection) {
		w := &schema.Work{}

		w.Type = strings.ToLower(s.Find(".issue-heading").Text())

		title := s.Find(".hlFld-Title")
		doi, ok := title.Find("a").Attr("href")
		if !ok {
			return
		}

		w.DOI = strings.TrimPrefix(doi, "/doi/")
		w.Title = util.CleanString(title.Text())
		w.Venue = s.Find(".epub-section__title").Text()
		w.Page = util.ParsePages(s.Find(".issue-item__detail").Find(".dot-separator").Text())

		s.Find(".rlist--inline").Each(func(i int, s *goquery.Selection) {
			if value, ok := s.Attr("aria-label"); !ok || value != "authors" {
				return
			}

			s.Find("a").Each(func(i int, s *goquery.Selection) {
				author, ok := s.Attr("title")
				if !ok {
					return
				}

				w.Authors = append(w.Authors, util.CleanString(author))
			})
		})

		date, ok := s.Find(".bookPubDate").Attr("data-title")
		if !ok {
			return
		}

		day, month, year, err := util.ParseDate(strings.TrimPrefix(date, "Published: "))
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Day = day
		w.Month = month
		w.Year = year

		works = append(works, w)
	})

	return works
}

// SearchACM searches the ACM digital library for a work
func SearchACM(ctx context.Context, w *schema.Work) ([]*schema.Work, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://dl.acm.org/action/doSearch?AllField="+url.QueryEscape(w.Title), nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return getAcmWorks(doc), nil
}
