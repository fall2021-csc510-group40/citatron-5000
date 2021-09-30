/*
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
*/
package util

import (
	"testing"
)

func expect(t *testing.T, value interface{}, expected interface{}) {
	if value == expected {
		return
	}

	t.Errorf("wanted %v got %v", expected, value)
}

func TestCleanString(t *testing.T) {
	expect(t, CleanString(" a  b  c "), "a b c")
}

func TestRemoveAllPunctuation(t *testing.T) {
	expect(t, RemoveAllPunctuation(";$3a/b-c;"), "3abc")
}

func TestGetMonth(t *testing.T) {
	m, err := GetMonth("july")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expect(t, m, 7)

	errInput := "abc"
	if _, err := GetMonth("abc"); err == nil {
		t.Errorf("expected error for month %v", errInput)
	}
}

func TestParseDate(t *testing.T) {
	d, m, y, err := ParseDate("Feb 3 2001")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expect(t, d, 3)
	expect(t, m, 2)
	expect(t, y, 2001)

	d, m, y, err = ParseDate("27 aug. 1997")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	expect(t, d, 27)
	expect(t, m, 8)
	expect(t, y, 1997)
}

func TestParsePages(t *testing.T) {
	expect(t, ParsePages("1-100"), "1-100")
	expect(t, ParsePages("1997 pp. 5-7"), "5-7")
	expect(t, ParsePages("pp. 9"), "9")
}
