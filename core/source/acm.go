package source

import (
	"core/schema"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

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

	doc.Find(".issue-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".hlFld-Title").Text()
		title = regexp.MustCompile(`\s\s+`).ReplaceAllString(title, " ")

		fmt.Println(title)
	})

	return nil, nil
}
