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
