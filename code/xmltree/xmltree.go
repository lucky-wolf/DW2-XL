package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// alternate/replacement library for golang's xml package
// this one doesn't use maps, but rather arrays, so that order is always preserved
// we're not trying to do everything, just a few things well, so ymmv

var ErrEOF = io.EOF

type XMLTree struct {
	Elements XMLValue // the one thing this cannot be is just a string, but an array of any of the others is allowed
}

type XMLValue struct {
	contents any // can be a string, comment, directive, processing instructions, element or an array of anything other than string
}

type XMLProcInst struct {
	xml.ProcInst
}

type XMLComment struct {
	xml.Comment
}

type XMLDirective struct {
	xml.Directive
}

type XMLElement struct {
	xml.StartElement
	XMLValue // can be a single string, or an array of child elements such as other elements or comments etc.
}

type Tokenizer interface {
	Token() (xml.Token, error)
	InputPos() (line, column int)
}

// returns an XMLTree by reading in the given file
func LoadFromFile(filename string) (tree XMLTree, err error) {

	// first we need a stream
	stream, err := os.Open(filename)
	if err != nil {
		return
	}
	defer stream.Close()

	// then we need to tokenize the stream
	err = tree.Read(stream)
	return
}

func (tree *XMLTree) Read(stream io.Reader) (err error) {

	// decode the stream into ourself
	err = tree.decode(xml.NewDecoder(stream))

	// eof is fine
	if err == io.EOF {
		err = nil
	}

	return
}

func (tree *XMLTree) decode(tokenizer Tokenizer) (err error) {
	// decode the root element(s)
	err = tree.Elements.decodeRoot(tokenizer)
	return
}

func (value *XMLValue) decodeRoot(tokenizer Tokenizer) (err error) {

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
			value.append(XMLComment{v.Copy()})
		case xml.Directive:
			value.append(XMLDirective{v.Copy()})
		case xml.ProcInst:
			value.append(XMLProcInst{v.Copy()})
		case xml.StartElement:
			if root != nil {
				err = SyntaxError(tokenizer, "only one root element", v)
				return
			}
			root = &XMLElement{StartElement: v.Copy()}
			err = root.decode(tokenizer)
			if err != nil {
				return
			}
			value.append(root)
		case xml.EndElement:
			err = SyntaxError(tokenizer, "anything else", v)
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

func (e *XMLElement) decode(tokenizer Tokenizer) (err error) {

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
			e.append(XMLComment{v.Copy()})
		case xml.Directive:
			e.append(XMLDirective{v.Copy()})
		case xml.ProcInst:
			e.append(XMLProcInst{v.Copy()})
		case xml.StartElement:
			child := &XMLElement{StartElement: v.Copy()}
			err = child.decode(tokenizer)
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
		}
	}
}

func FQN(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	return name.Space + ":" + name.Local
}

func SyntaxError(tokenizer Tokenizer, expected, found any) error {
	line, _ := tokenizer.InputPos()
	return &xml.SyntaxError{Msg: fmt.Sprintf("expected: %v, found: %v", expected, found), Line: line}
}

func (e XMLProcInst) String() string {

	sb := strings.Builder{}
	sb.WriteString("<?")
	sb.WriteString(e.Target)
	if len(e.Inst) > 0 {
		sb.WriteByte(' ')
		sb.Write(e.Inst)
	}
	sb.WriteString("?>")

	return sb.String()
}

func (e XMLDirective) String() string {

	sb := strings.Builder{}
	sb.WriteString("<!")
	sb.Write(e.Directive)
	sb.WriteString(">")

	return sb.String()
}

func (e XMLComment) String() string {

	sb := strings.Builder{}
	sb.WriteString("<!--")
	sb.Write(e.Comment)
	sb.WriteString("-->")

	return sb.String()
}

func (e XMLElement) String() string {

	sb := new(strings.Builder)

	// write the start token with attributes
	sb.WriteByte('<')
	sb.WriteString(FQN(e.Name))
	for i := range e.Attr {
		sb.WriteByte(' ')
		sb.WriteString(FQN(e.Attr[i].Name))
		sb.WriteByte('=')
		sb.WriteByte('"')
		//todo: we need to map illegal chars to their &xx; equivalents
		sb.WriteString(e.Attr[i].Value)
		sb.WriteByte('"')
	}

	if e.XMLValue.Empty() {
		sb.WriteString(" />")
	} else {
		sb.WriteByte('>')

		// ask XMLItem to express itself
		sb.WriteString(e.XMLValue.String())

		// write the closure
		sb.WriteString("</")
		sb.WriteString(e.Name.Local)
		sb.WriteByte('>')
	}

	// sb.WriteByte('\n')

	return sb.String()
}

func (e XMLValue) Empty() bool {

	// absolute empty
	if e.contents == nil {
		return true
	}

	// simple string contents
	if v, ok := e.contents.(string); ok {
		return len(strings.TrimSpace(v)) == 0
	}

	return false
}

func (e XMLValue) String() string {

	// return the empty string if we're empty
	if e.Empty() {
		return ""
	}

	// simple case is just string contents
	if v, ok := e.contents.(string); ok {
		return v
	}

	// everything else is one or more child objects
	sb := new(strings.Builder)
	switch v := e.contents.(type) {
	case XMLComment:
		sb.WriteByte('\n')
		sb.WriteString(v.String())
	case XMLDirective:
		sb.WriteByte('\n')
		sb.WriteString(v.String())
	case XMLProcInst:
		sb.WriteByte('\n')
		sb.WriteString(v.String())
	case XMLElement:
		sb.WriteByte('\n')
		sb.WriteString(v.String())
	case []any:
		for i := range v {
			sb.WriteByte('\n')
			sb.WriteString(fmt.Sprint(v[i]))
		}
	}
	return sb.String()
}

func (e XMLTree) String() string {
	return e.Elements.String()
}
