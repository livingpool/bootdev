package util

import (
	"reflect"
	"slices"
	"testing"
)

func TestSortPages(t *testing.T) {
	tests := []struct {
		name     string
		input    []Page
		expected []Page
	}{
		{
			name: "order count descending",
			input: []Page{
				{URL: "url1", Count: 5},
				{URL: "url2", Count: 1},
				{URL: "url3", Count: 3},
				{URL: "url4", Count: 10},
				{URL: "url5", Count: 7},
			},
			expected: []Page{
				{URL: "url4", Count: 10},
				{URL: "url5", Count: 7},
				{URL: "url1", Count: 5},
				{URL: "url3", Count: 3},
				{URL: "url2", Count: 1},
			},
		},
		{
			name: "alphabetize",
			input: []Page{
				{URL: "d", Count: 1},
				{URL: "a", Count: 1},
				{URL: "e", Count: 1},
				{URL: "b", Count: 1},
				{URL: "c", Count: 1},
			},
			expected: []Page{
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
				{URL: "c", Count: 1},
				{URL: "d", Count: 1},
				{URL: "e", Count: 1},
			},
		},
		{
			name: "order count then alphabetize",
			input: []Page{
				{URL: "d", Count: 2},
				{URL: "a", Count: 1},
				{URL: "e", Count: 3},
				{URL: "b", Count: 1},
				{URL: "c", Count: 2},
			},
			expected: []Page{
				{URL: "e", Count: 3},
				{URL: "c", Count: 2},
				{URL: "d", Count: 2},
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
			},
		},
		{
			name:     "empty list",
			input:    []Page{},
			expected: []Page{},
		},
		{
			name:     "nil list",
			input:    nil,
			expected: nil,
		},
		{
			name: "one key",
			input: []Page{
				{URL: "url1", Count: 1},
			},
			expected: []Page{
				{URL: "url1", Count: 1},
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slices.SortFunc(tc.input, SortPages)
			if !reflect.DeepEqual(tc.input, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, tc.input)
			}
		})
	}
}
