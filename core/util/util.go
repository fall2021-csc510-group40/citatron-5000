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

// Package util provides several functions for convenient string-manipulation as well as other convenient constants, types and functions.
package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var duplicateSpaceRegex = regexp.MustCompile(`\s\s+`)
var punctuationRegex = regexp.MustCompile(`[^a-zA-Z\d]`)

var dayRegex = regexp.MustCompile(`\b\d{1,2}\b`)
var monthRegex = regexp.MustCompile(`\b[a-zA-Z]+\b`)
var yearRegex = regexp.MustCompile(`\b\d{4}\b`)

var pageRegex = regexp.MustCompile(`\d+-\d+`)
var numberRegex = regexp.MustCompile(`\d+`)

// CleanString removes all unnecessary whitespace
func CleanString(s string) string {
	return duplicateSpaceRegex.ReplaceAllString(strings.TrimSpace(s), " ")
}

// RemoveAllPunctuation removes all non-alphanumeric characters
func RemoveAllPunctuation(s string) string {
	return punctuationRegex.ReplaceAllString(s, "")
}

// GetMonth returns the month number from string
func GetMonth(m string) (int, error) {
	switch strings.ToLower(m) {
	case "jan", "january":
		return 1, nil
	case "feb", "february":
		return 2, nil
	case "mar", "march":
		return 3, nil
	case "apr", "april":
		return 4, nil
	case "may":
		return 5, nil
	case "jun", "june":
		return 6, nil
	case "jul", "july":
		return 7, nil
	case "aug", "august":
		return 8, nil
	case "sept", "september":
		return 9, nil
	case "oct", "october":
		return 10, nil
	case "nov", "november":
		return 11, nil
	case "dec", "december":
		return 12, nil
	}

	return 0, errors.New("invalid month")
}

// ParseDate parses a month, day, year string into its parts
func ParseDate(d string) (int, int, int, error) {
	day, err := strconv.Atoi(dayRegex.FindString(d))
	if err != nil || day < 0 || day > 31 {
		return 0, 0, 0, errors.New("bad day")
	}

	month, err := GetMonth(monthRegex.FindString(d))
	if err != nil {
		return 0, 0, 0, errors.New("bad month")
	}

	year, err := strconv.Atoi(yearRegex.FindString(d))
	if err != nil || year < 0 {
		return 0, 0, 0, errors.New("bad year")
	}

	return day, month, year, nil
}

// ParsePages parses a page range string into its simplest form
func ParsePages(s string) string {
	s = strings.ReplaceAll(s, "–", "-")

	match := pageRegex.FindString(s)
	if match != "" {
		return match
	}

	split := strings.Split(s, "pp")
	return numberRegex.FindString(split[len(split)-1])
}
