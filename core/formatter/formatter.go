package formatter

import (
	"core/schema"
	"strconv"
	"strings"
)

func BibtexFormat(paper *schema.Work) string {
	var citation string = "@article{ \n" //need to add the key

	if len(paper.Authors) > 0 {
		citation += "Author= { "
		for i := 0; i < len(paper.Authors)-1; i++ {
			citation += paper.Authors[i] + ", "
		}
		citation += paper.Authors[len(paper.Authors)-1]
		citation += " },\n"
	}

	if len(paper.Title) > 0 {
		citation += "Title= { " + paper.Title + " },\n"
	}

	if len(paper.DOI) > 0 {
		citation += "DOI= { " + paper.DOI + " },\n"
	}

	if len(paper.Arxiv) > 0 {
		citation += "ARXIV= { " + paper.Arxiv + " },\n"
	}
	if len(paper.ISBN) > 0 {
		citation += "ISBN= { " + paper.Arxiv + " },\n"
	}

	if len(paper.Venue) > 0 {
		citation += "Journal= { " + paper.Venue + " },\n"
	}
	if paper.Month != 0 {
		citation += "Month= { " + strconv.Itoa(paper.Month) + " }, \n"
	}

	if paper.Year != 0 {
		citation += "Year= { " + strconv.Itoa(paper.Year) + " }, \n"
	}

	if paper.Page != "" {
		citation += "Page= { " + paper.Page + " }, \n"
	}
	citation += "},\n"
	return citation
}

func PlaintextFormat(paper *schema.Work) string {
	var citation string = ""
	if len(paper.Title) > 0 {
		citation += paper.Title + ", "
	}
	if len(paper.Authors) > 0 {
		for _, author := range paper.Authors {
			citation += author + ", "
		}
	}

	if len(paper.Venue) > 0 {
		citation += paper.Venue + ", "
	}
	if paper.Month != 0 {
		citation += strconv.Itoa(paper.Month) + ", "
	}
	if paper.Year != 0 {
		citation += strconv.Itoa(paper.Year) + ","
	}
	citation = strings.TrimSpace(citation)
	citation = strings.TrimSuffix(citation, ",")
	citation += "."
	return citation
}
