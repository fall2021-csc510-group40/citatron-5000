package formatter

import (
	"core/schema"
	"strconv"
	"strings"
	"time"
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
		citation += "Month= { " + strconv.Itoa(paper.Year) + " }, \n"
	}

	if paper.Year != 0 {
		citation += "Year= { " + strconv.Itoa(paper.Year) + " }, \n"
	}

	if paper.Page != "" {
		citation += "Page= { " + paper.Page + " }, \n"
	}

	// need to check this logic
	// if len(paper.Authors) > 0 {
	// 	citation += "author={"
	// 	for i, author := range paper.Authors {
	// 		parts := strings.Split(author, " ")
	// 		if len(parts) > 1 {
	// 			citation += parts[1] + ", " + parts[0]
	// 		} else {
	// 			citation += parts[0]
	// 		}
	// 		if i != (len(paper.Authors) - 1) {
	// 			citation += "and"
	// 		}
	// 	}
	citation += "},\n"
	// }
	return citation
}

func PlaintextFormat(paper *schema.Work) string {
	var citation string = ""
	if len(paper.Authors) > 0 {
		for i, author := range paper.Authors {
			parts := strings.Split(author, " ")
			if len(parts) > 1 {
				citation += string(parts[0][0]) + ". " + parts[1]
			} else {
				citation += parts[0]
			}
			if i == (len(paper.Authors) - 2) {
				citation += "and"
			}
		}
		citation += ", "
	}
	if len(paper.Title) > 0 {
		citation += "\"" + paper.Title + ",\""
	}
	if len(paper.Venue) > 0 {
		citation += paper.Venue + ", "
	}
	if paper.Month != 0 {
		citation += time.Month(paper.Month).String() + ", "
	}
	if paper.Year != 0 {
		citation += strconv.Itoa(paper.Year)
	}
	return citation
}
