package formatter

import (
	"core/schema"
	"strconv"
	"strings"
)

// BibtexFormat formats a work in Bibtex style
func BibtexFormat(work *schema.Work) string {
	citation := "@article{ YOUR_KEY_HERE,\n"

	if len(work.Authors) > 0 {
		citation += "Author= { "
		for i := 0; i < len(work.Authors)-1; i++ {
			citation += work.Authors[i] + ", "
		}
		citation += work.Authors[len(work.Authors)-1]
		citation += " },\n"
	}

	if len(work.Title) > 0 {
		citation += "Title= { " + work.Title + " },\n"
	}

	if len(work.DOI) > 0 {
		citation += "DOI= { " + work.DOI + " },\n"
	}

	if len(work.Arxiv) > 0 {
		citation += "ARXIV= { " + work.Arxiv + " },\n"
	}

	if len(work.ISBN) > 0 {
		citation += "ISBN= { " + work.Arxiv + " },\n"
	}

	if len(work.Venue) > 0 {
		citation += "Journal= { " + work.Venue + " },\n"
	}
	if work.Month != 0 {
		citation += "Month= { " + strconv.Itoa(work.Month) + " }, \n"
	}

	if work.Year != 0 {
		citation += "Year= { " + strconv.Itoa(work.Year) + " }, \n"
	}

	if work.Page != "" {
		citation += "Page= { " + work.Page + " }, \n"
	}

	citation += "},\n"
	return citation
}

// PlaintextFormat formats a work in a simple plaintext format
func PlaintextFormat(work *schema.Work) string {
	citation := ""

	if len(work.Title) > 0 {
		citation += work.Title + ", "
	}

	if len(work.Authors) > 0 {
		for _, author := range work.Authors {
			citation += author + ", "
		}
	}

	if len(work.Venue) > 0 {
		citation += work.Venue + ", "
	}

	if work.Month != 0 {
		citation += strconv.Itoa(work.Month) + ", "
	}

	if work.Year != 0 {
		citation += strconv.Itoa(work.Year) + ","
	}

	citation = strings.TrimSpace(citation)
	citation = strings.TrimSuffix(citation, ",")
	citation += "."

	return citation
}
