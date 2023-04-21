package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// interface we require for our tokenization
type Tokenizer interface {
	Token() (xml.Token, error)
	InputPos() (line, column int)
}

// errors are usually, but not always, output this way
func SyntaxError(tokenizer Tokenizer, expected, found any) error {
	line, _ := tokenizer.InputPos()
	return &xml.SyntaxError{Msg: fmt.Sprintf("expected: %v, found: %v", expected, found), Line: line}
}

// returns an XMLTree by reading in the given file
func LoadFromFile(filename string) (tree *XMLTree, err error) {

	// first we need a stream
	stream, err := os.Open(filename)
	if err != nil {
		return
	}
	defer stream.Close()

	// then we need to tokenize the stream
	tree = new(XMLTree)
	err = tree.Read(stream)
	return
}

// returns an XMLTree by reading from the stream
func (tree *XMLTree) Read(stream io.Reader) (err error) {

	// decode the stream into ourself
	err = tree.Decode(xml.NewDecoder(stream))

	// eof is fine
	if err == io.EOF {
		err = nil
	}

	return
}

// you can supply your own tokenizer if desired
func (tree *XMLTree) Decode(tokenizer Tokenizer) (err error) {
	err = tree.Elements.DecodeRoot(tokenizer)
	return
}

// note: from this point forward it seems unlikely that you'd be calling into the library at this level
// note: but we're leaving it possible because "why not?"

// decodes our root, which can only contain one element, but may contain any number of comments, prodinst, etc.
// we don't care about such things, but we do our best to faithfully capture them
func (value *XMLValue) DecodeRoot(tokenizer Tokenizer) (err error) {

	var root *XMLElement

	for {

		var token xml.Token
		token, err = tokenizer.Token()
		if err != nil {
			return
		}

		switch v := token.(type) {
		case xml.CharData:
			// at the root, we just ignore char data which should be whitespace
			s := strings.TrimSpace(string(v))
			if len(s) != 0 {
				err = SyntaxError(tokenizer, "whitespace only", s)
				return
			}
		case xml.Comment:
			value.append(&XMLComment{v.Copy()})
		case xml.Directive:
			value.append(&XMLDirective{v.Copy()})
		case xml.ProcInst:
			value.append(&XMLProcInst{v.Copy()})
		case xml.StartElement:
			if root != nil {
				err = SyntaxError(tokenizer, "only one root element", v)
				return
			}
			root = &XMLElement{StartElement: v.Copy()}
			err = root.Decode(tokenizer)
			if err != nil {
				return
			}
			value.append(root)
		case xml.EndElement:
			err = SyntaxError(tokenizer, "anything else", v)
			return
		default:
			err = UnknownEntity(token)
			return
		}
	}
}

func (value *XMLValue) append(item any) {

	// no value, then just store this one item
	if value.contents == nil {
		value.contents = item
		return
	}

	// if we don't already have multiple children, convert us to an array
	children, ok := value.contents.([]any)
	if !ok {
		children = []any{value.contents}
	}

	// append us
	value.contents = append(children, item)
	return
}

func (e *XMLElement) Decode(tokenizer Tokenizer) (err error) {

	var text string

	for {

		var token xml.Token
		token, err = tokenizer.Token()
		if err != nil {
			return
		}

		switch v := token.(type) {
		case xml.CharData:
			// at the element level, strings might be contents
			text += string(v)
		case xml.Comment:
			e.append(&XMLComment{v.Copy()})
		case xml.Directive:
			e.append(&XMLDirective{v.Copy()})
		case xml.ProcInst:
			e.append(&XMLProcInst{v.Copy()})
		case xml.StartElement:
			child := &XMLElement{StartElement: v.Copy()}
			err = child.Decode(tokenizer)
			if err != nil {
				return
			}
			e.append(child)
		case xml.EndElement:
			if v.Name.Local != e.Name.Local {
				err = SyntaxError(tokenizer, e.Name.Local, v.Name.Local)
			} else if e.contents == nil {
				e.contents = text
			} else if len(strings.TrimSpace(text)) != 0 {
				err = SyntaxError(tokenizer, "cannot combine text and children elements", text)
				return
			}
			return
		default:
			err = UnknownEntity(token)
			return
		}
	}
}
