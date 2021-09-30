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
package schema

import "testing"

func expect(t *testing.T, value interface{}, expected interface{}) {
	if value == expected {
		return
	}

	t.Errorf("wanted %v got %v", expected, value)
}

func TestWorkNormalize(t *testing.T) {
	w := &Work{}
	w.Title = "  Some   Title  "
	w.Authors = append(w.Authors, "  James   Smith  ")

	if err := w.Normalize(); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expect(t, w.Title, "Some Title")
	expect(t, w.Authors[0], "James Smith")
}

func TestWorkCoalesce(t *testing.T) {
	w0 := &Work{}
	w0.Title = "  some   Title  "
	w0.Authors = append(w0.Authors, "  James-Smith. ")
	w0.Arxiv = "arxiv"

	if err := w0.Normalize(); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	w1 := &Work{}
	w1.Title = "  Some title"
	w1.Authors = append(w1.Authors, "james smith")
	w1.DOI = "doi"

	if err := w1.Normalize(); err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expect(t, w0.Hash, w1.Hash)

	w0.Coalesce(w1)

	expect(t, w0.Arxiv, "arxiv")
	expect(t, w0.DOI, "doi")
}
