package formatter

import (
	"core/schema"
	"strconv"
	"strings"
)

func BibtexFormat(paper *schema.Work) string {
	var citation string = "@article{"
	for _, author := range paper.Authors {
		citation += author + "_"
	}
	citation += strconv.Itoa(paper.Year) + ", title={"
	citation += paper.Title + "}, DOI={"
	citation += paper.DOI + "}, journal={"
	citation += paper.Venue + "}, author={"
	for i, author := range paper.Authors {
		parts := strings.Split(author, " ")
		if len(parts) > 0 {
			citation += parts[1] + ", " + parts[0]
		} else {
			citation += parts[0]
		}
		if i != (len(paper.Authors) - 1) {
			citation += "and"
		}
	}
	citation += "}, year={" + strconv.Itoa(paper.Year) + "}}"

	return citation
}

func PlaintextFormat(paper *schema.Work) string {
	var citation string = ""

	return citation
}
