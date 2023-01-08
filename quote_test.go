package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func FixQuotedQuote(line []byte) []byte {
	recordBuffer := bytes.Split(line, []byte(","))
	for i, field := range recordBuffer {
		quoteStart := bytes.IndexByte(field, '"')
		if quoteStart < 0 {
			continue
		}

		quoteEnd := bytes.LastIndexByte(field, '"')

		recordBuffer[i] = append(field[:quoteStart+1], bytes.ReplaceAll(field[quoteStart+1:quoteEnd], []byte("\""), []byte("\"\""))...)
		recordBuffer[i] = append(recordBuffer[i], field[quoteEnd:]...)

	}

	return bytes.Join(recordBuffer, []byte(","))
}

func AddQuotedQuote(line []byte) []byte {
	recordBuffer := bytes.Split(line, []byte(","))
	for i, field := range recordBuffer {
		quoteStart := bytes.IndexByte(field, '"')
		if quoteStart < 0 {
			continue
		}

		recordBuffer[i] = append([]byte("\""), bytes.ReplaceAll(field, []byte("\""), []byte("\"\""))...)
		recordBuffer[i] = append(recordBuffer[i], '"')

	}

	return bytes.Join(recordBuffer, []byte(","))
}

var badCSV = []readTest{
	{
		Name:    "https://github.com/golang/go/issues/56903",
		Input:   `"field1","field2"","field3"`,
		Output:  [][]string{{"field1", "field2\"", "field3"}},
		Formats: []Format{FixQuotedQuote},
	},
	{
		Name: "https://github.com/golang/go/issues/56329 example 1",
		Input: `
c0,c1,c2
abc,123,
"abc",123,
"a"b"c",123,
""ab"c",123,
""abc"",123,
"a"b"c",123,
`,
		Output:  [][]string{{"c0", "c1", "c2"}, {"abc", "123", ""}, {"abc", "123", ""}, {"a\"b\"c", "123", ""}, {"\"ab\"c", "123", ""}, {"\"abc\"", "123", ""}, {"a\"b\"c", "123", ""}},
		Formats: []Format{FixQuotedQuote},
	},
	{
		Name: "https://github.com/golang/go/issues/56329 example 2",
		Input: `
c0,c1,c2
abc,123,
"a"bc"",123,
a"b"c,123,
`,
		Output:  [][]string{{"c0", "c1", "c2"}, {"abc", "123", ""}, {"a\"bc\"", "123", ""}, {"a\"b\"c", "123", ""}},
		Formats: []Format{FixQuotedQuote},
	},
	{
		Name:    "https://github.com/golang/go/issues/55069",
		Input:   `"a"b,c`,
		Output:  [][]string{{"\"a\"b", "c"}},
		Formats: []Format{AddQuotedQuote},
	},
}

func TestLazyQuote(t *testing.T) {
	for _, tt := range badCSV {
		t.Run(tt.Name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.Input))
			r.LazyQuotes = true
			out, err := r.ReadAll(tt.Formats...)
			if err != nil {
				t.Fatalf("unexpected Readall() error: %v", err)
			}
			if !reflect.DeepEqual(out, tt.Output) {
				t.Fatalf("ReadAll() output:\ngot  %q\nwant %q", out, tt.Output)
			}
		})
	}

}
