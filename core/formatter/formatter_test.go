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
	"fmt"
	"testing"
)

func expect(t *testing.T, value interface{}, expected interface{}) {
	if value == expected {
		return
	}

	t.Errorf("wanted %v got %v", expected, value)
}

func TestFormatter(t *testing.T) {
	w := &schema.Work{}
	w.Title = "Interesting   Title"
	w.Authors = append(w.Authors, "James   Smith,Rose Walnut")
	w.DOI = "10.0.0.xxx"
	w.Page = "23-34"
	w.Year = 2021
	w.Month = 11
	w.Venue = "Research journal"

	resBibtex := BibtexFormat(w)
	expectedResBibtex := "@article{ YOUR_KEY_HERE,\n    Author= { James   Smith,Rose Walnut },\n    Title= { Interesting   Title },\n    DOI= { 10.0.0.xxx },\n    Journal= { Research journal },\n    Month= { 11 },\n    Year= { 2021 },\n    Page= { 23-34 },\n},\n"
	expect(t, resBibtex, expectedResBibtex)

	resPlain := PlaintextFormat(w)
	expectedResPlain := "Interesting   Title, James   Smith,Rose Walnut, Research journal, 11, 2021."
	fmt.Println(resPlain, expectedResPlain)
	expect(t, resPlain, expectedResPlain)
}
