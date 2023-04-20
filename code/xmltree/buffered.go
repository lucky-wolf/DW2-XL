package xmltree

// some lower level parsing functions

import (
	"encoding/xml"
	"io"
)

type BufferedTokenizer struct {
	tokenizer *xml.Decoder
	peek      xml.Token
}

func NewTokenizer(stream io.Reader) (t *BufferedTokenizer) {

	t = &BufferedTokenizer{tokenizer: xml.NewDecoder(stream)}

	return
}

func (t *BufferedTokenizer) Peek() (token xml.Token, err error) {

	// if we don't yet have a peek, grab it
	if t.peek == nil {
		t.peek, err = t.tokenizer.Token()
		if err != nil {
			return
		}
	}

	// return the peek token
	token = t.peek
	return
}

func (t *BufferedTokenizer) Token() (token xml.Token, err error) {

	// grab buffered token if it exists (and remove it from peek)
	if t.peek != nil {
		token = t.peek
		t.peek = nil
		return
	}

	// or get a new token (peek is nil)
	token, err = t.tokenizer.Token()
	return
}

// advances past the current peek token (if any) - returns true if there was one
func (t *BufferedTokenizer) Advance() bool {
	if t.peek == nil {
		return false
	}
	t.peek = nil
	return true
}
