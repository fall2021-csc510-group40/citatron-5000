package formatter

import (
	"core/schema"
	"strconv"
	"strings"
	"time"
)

func BibtexFormat(paper *schema.Work) string {
	var citation string = "@article{"
	if len(paper.Authors) > 0 {
		for _, author := range paper.Authors {
			citation += author + "_"
		}
	}
	if paper.Year != 0 {
		citation += strconv.Itoa(paper.Year)
	}
	if paper.Year == 0 && len(paper.Authors) == 0 {
		citation += "empty"
	}
	if len(paper.Title) > 0 {
		citation += ", title={" + paper.Title + "}, "
	}
	if len(paper.DOI) > 0 {
		citation += "DOI={" + paper.DOI + "}, "
	}
	if len(paper.Venue) > 0 {
		citation += "journal={" + paper.Venue + "}, "
	}
	if len(paper.Authors) > 0 {
		citation += "author={"
		for i, author := range paper.Authors {
			parts := strings.Split(author, " ")
			if len(parts) > 1 {
				citation += parts[1] + ", " + parts[0]
			} else {
				citation += parts[0]
			}
			if i != (len(paper.Authors) - 1) {
				citation += "and"
			}
		}
		citation += "}, "
	}

	if paper.Year != 0 {
		citation += "year={" + strconv.Itoa(paper.Year) + "}}"
	}
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
