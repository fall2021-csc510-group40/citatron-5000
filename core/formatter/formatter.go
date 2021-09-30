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
		citation += "    Author= { "
		for i := 0; i < len(work.Authors)-1; i++ {
			citation += work.Authors[i] + ", "
		}
		citation += work.Authors[len(work.Authors)-1]
		citation += " },\n"
	}

	if len(work.Title) > 0 {
		citation += "    Title= { " + work.Title + " },\n"
	}

	if len(work.DOI) > 0 {
		citation += "    DOI= { " + work.DOI + " },\n"
	}

	if len(work.Arxiv) > 0 {
		citation += "    ARXIV= { " + work.Arxiv + " },\n"
	}

	if len(work.ISBN) > 0 {
		citation += "    ISBN= { " + work.Arxiv + " },\n"
	}

	if len(work.Venue) > 0 {
		citation += "    Journal= { " + work.Venue + " },\n"
	}
	if work.Month != 0 {
		citation += "    Month= { " + strconv.Itoa(work.Month) + " }, \n"
	}

	if work.Year != 0 {
		citation += "    Year= { " + strconv.Itoa(work.Year) + " }, \n"
	}

	if work.Page != "" {
		citation += "    Page= { " + work.Page + " }, \n"
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
