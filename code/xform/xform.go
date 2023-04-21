package xform

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// relatively simple csv reader that uses the format:
// search=<token>,...		tells us which tokens to search on
// replace=<token>,...      tells us which tokens are replacement ones
//                          note: any columns not listed will simply be ignored
// <token>,<token>,<token>,...      defines what the follow columnar data corresponds to
// value, value, value, ...			the values that go with the corresponding tokens
//                                  note: any values that are missing are assumed to be same as previous line
//							please see "XL Armor.csv" for an example file

func LoadFromFile(filename string) (xs XScript, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	err = xs.Decode(file)

	return
}

type XScript struct {
	Search  []string   // the tags to search for
	Replace []string   // the tags to replace on
	Header  []string   // the tags that define the data that follows
	Values  [][]string // the search or replacement values
}

func StartsWith(s, sub string) bool {
	return len(s) >= len(sub) && s[:len(sub)] == sub
}

type LineScanner struct {
	*bufio.Scanner
	trimSpace      bool
	skipBlankLines bool
	text           string
	line           int
}

func NewLineScanner(r io.Reader, trimSpace bool, skipBlankLines bool) *LineScanner {
	return &LineScanner{Scanner: bufio.NewScanner(r), trimSpace: trimSpace}
}

func (scanner *LineScanner) scan() bool {
	if !scanner.Scanner.Scan() {
		return false
	}
	scanner.line++
	scanner.text = scanner.Scanner.Text()
	if scanner.trimSpace {
		scanner.text = strings.TrimSpace(scanner.text)
	}
	return true
}

func (scanner *LineScanner) Text() string {
	return scanner.text
}

func (scanner *LineScanner) Scan() (scanned bool) {
	if !scanner.skipBlankLines {
		return scanner.scan()
	}
	for scanner.scan() {
		if len(scanner.text) != 0 {
			return true
		}
	}
	return false
}

func (scanner *LineScanner) Line() int {
	return scanner.line
}

var regexTag = regexp.MustCompile(`<(\w+)>`)
var ErrMismatchColumnLength = errors.New("mismatched column length")

func (xs *XScript) Decode(r io.Reader) (err error) {

	// process our file line by line
	scanner := NewLineScanner(r, true, true)

	// automagically add line number on error
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s at line: %d", err, scanner.Line())
		}
	}()

	// extract the search tags
	xs.Search, err = decodeDeclaration(scanner, "search=")
	if err != nil {
		return
	}

	// extract the search tags
	xs.Replace, err = decodeDeclaration(scanner, "replace=")
	if err != nil {
		return
	}

	// for each line, process instructions
	for scanner.Scan() {

		// break this line into columns
		columns := strings.Split(scanner.Text(), ",")
		if regexTag.MatchString(columns[0]) {

			// this is the column header
			if xs.Header == nil {
				xs.Header = columns
				continue
			}

			// all redundant headers must match
			if len(xs.Header) != len(columns) {
				err = ErrMismatchColumnLength
				return
			}
			for i := range xs.Header {
				if xs.Header[i] != columns[i] {
					err = fmt.Errorf("Mismatched column definitions")
					return
				}
			}
			continue
		}

		// require all rows have the same count of columns
		if len(columns) != len(xs.Header) {
			err = fmt.Errorf("column length cannot vary between %d and %d", len(columns), len(xs.Header))
			return
		}

		// a values row
		xs.Values = append(xs.Values, columns)
	}

	return
}

func decodeDeclaration(scanner *LineScanner, key string) (tags []string, err error) {

	if !scanner.Scan() || !StartsWith(scanner.Text(), key) {
		err = errors.New("%s declaration not found")
		return
	}

	// extract the search tags
	tags = strings.Split(scanner.Text()[len(key):], ",")
	return
}
