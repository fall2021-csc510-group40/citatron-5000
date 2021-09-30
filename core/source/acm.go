/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package source

import (
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

// SourceSearchACM searches the ACM digital library for a work
func SourceSearchACM(w *schema.Work) ([]*schema.Work, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Get("https://dl.acm.org/action/doSearch?AllField=" + url.QueryEscape(w.Title))
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
