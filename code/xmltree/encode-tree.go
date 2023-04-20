package xmltree

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// writes ourself out to the given file (file is created / truncated to just our contents)
func (tree *XMLTree) WriteToFile(filename string) (err error) {

	// first we need a stream
	stream, err := os.Create(filename)
	if err != nil {
		return
	}
	defer stream.Close()

	encoder := NewEncoder(stream)
	defer encoder.Close()

	// then we need to write to the stream
	err = tree.Encode(encoder)

	return
}

// encode ourself into stream using default encoding settings (no formatting)
func (tree *XMLTree) Write(stream io.Writer) (err error) {
	err = tree.Encode(NewEncoder(stream))
	return
}

func (tree *XMLTree) Encode(encoder FormattedEncoder) (err error) {
	err = encoder.Indent(false, 0, true)
	if err != nil {
		return
	}
	err = tree.Elements.Encode(encoder)
	if err != nil {
		return
	}

	err = encoder.Indent(true, 0, false)
	if err != nil {
		return
	}

	err = encoder.Close()
	return
}

func (e *XMLValue) Encode(encoder FormattedEncoder) (err error) {

	// return the empty string if we're empty
	if e.Empty() {
		return
	}

	// simple case is just string contents
	if v, ok := e.contents.(string); ok {
		err = WriteEscapedText(v, encoder, false)
		return
	}

	// everything else is one or more child objects
	err = encodeChildren(e.contents, encoder)

	return
}

// this relies upon intimate knowledge of implementation of the xml tree
func encodeChildren(e any, encoder FormattedEncoder) (err error) {

	switch v := e.(type) {
	case *XMLComment:
		err = v.Encode(encoder)
	case *XMLDirective:
		err = v.Encode(encoder)
	case *XMLProcInst:
		err = v.Encode(encoder)
	case *XMLElement:
		err = v.Encode(encoder)
	case []any:
		for i, e := range v {
			if i != 0 {
				err = encoder.Indent(true, 0, true)
				if err != nil {
					return
				}
			}
			switch v := e.(type) {
			case *XMLComment:
				err = v.Encode(encoder)
			case *XMLDirective:
				err = v.Encode(encoder)
			case *XMLProcInst:
				err = v.Encode(encoder)
			case *XMLElement:
				err = v.Encode(encoder)
			default:
				err = UnknownEntity(e)
			}
			if err != nil {
				return
			}
		}
	default:
		err = UnknownEntity(e)
	}
	return
}

func (e *XMLElement) Encode(encoder FormattedEncoder) (err error) {

	// write the start token with attributes
	_, err = encoder.WriteString("<" + e.Name.Local)
	if err != nil {
		return
	}

	if e.Name.Space != "" {
		a := xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: e.Name.Space}
		err = encoder.WriteByte(' ')
		if err != nil {
			return
		}
		err = EncodeAttr(a, encoder)
		if err != nil {
			return
		}
	}

	for i := range e.Attr {
		err = encoder.WriteByte(' ')
		if err != nil {
			return
		}
		err = EncodeAttr(e.Attr[i], encoder)
		if err != nil {
			return
		}
	}

	if e.XMLValue.Empty() {
		_, err = encoder.WriteString(" />")
	} else {
		// finish the start tag
		err = encoder.WriteByte('>')
		if err != nil {
			return
		}

		// for anything but a string, we need to start indenting deeper
		if !e.XMLValue.IsSimple() {
			err = encoder.Indent(true, 1, true)
			if err != nil {
				return
			}
		}

		// ask XMLValue to express itself
		err = e.XMLValue.Encode(encoder)
		if err != nil {
			return
		}

		// for anything but a string, we need to pop out a level and put our end tag there
		if !e.XMLValue.IsSimple() {
			err = encoder.Indent(true, -1, true)
			if err != nil {
				return
			}
		}

		// write the end tag
		_, err = encoder.WriteString("</" + e.Name.Local + ">")
	}

	return
}

func EncodeAttr(a xml.Attr, encoder FormattedEncoder) (err error) {
	if a.Name.Space != "" {
		_, err = encoder.WriteString(a.Name.Space + ":")
		if err != nil {
			return
		}
	}

	_, err = encoder.WriteString(a.Name.Local + "=")
	if err != nil {
		return
	}
	_, err = encoder.WriteString(`"` + EscapeString(a.Value) + `"`)
	return
}

func QuotedString(value string) string {
	return `"` + EscapeString(value) + `"`
}

var (
	escQuote = []byte("&#34;") // shorter than "&quot;"
	escTick  = []byte("&#39;") // shorter than "&apos;"
	escAmp   = []byte("&amp;")
	escLT    = []byte("&lt;")
	escGT    = []byte("&gt;")
	escTab   = []byte("&#x9;")
	escNL    = []byte("&#xA;")
	escCR    = []byte("&#xD;")
	escFF    = []byte("\uFFFD") // Unicode replacement character
)

// EscapeString returns the properly escaped XML equivalent of the plain text data s
func EscapeString(s string) string {
	sb := &strings.Builder{}
	WriteEscapedText(s, sb, true)
	return sb.String()
}

func WriteEscapedText(s string, sb ByteAndStringWriter, newlines bool) (err error) {
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
		case '"':
			esc = escQuote
		case '\'':
			esc = escTick
		case '&':
			esc = escAmp
		case '<':
			esc = escLT
		case '>':
			esc = escGT
		case '\t':
			if newlines {
				esc = escTab
			}
		case '\n':
			if newlines {
				esc = escNL
			}
		case '\r':
			if newlines {
				esc = escCR
			}
		default:
			if !IsInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = escFF
				break
			}
			continue
		}
		_, err = sb.WriteString(s[last : i-width])
		if err != nil {
			return
		}
		_, err = sb.Write(esc)
		if err != nil {
			return
		}
		last = i
	}
	_, err = sb.WriteString(s[last:])

	return
}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of https://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func IsInCharacterRange(r rune) bool {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

func (e *XMLComment) Encode(w ByteAndStringWriter) (err error) {
	_, err = w.WriteString("<!--")
	if err != nil {
		return
	}
	_, err = w.Write(e.Comment)
	if err != nil {
		return
	}
	_, err = w.WriteString("-->")
	return
}

func (e *XMLProcInst) Encode(w ByteAndStringWriter) (err error) {
	_, err = w.WriteString("<?")
	if err != nil {
		return
	}
	_, err = w.WriteString(e.Target)
	if err != nil {
		return
	}
	if len(e.Inst) > 0 {
		err = w.WriteByte(' ')
		if err != nil {
			return
		}
		_, err = w.Write(e.Inst)
		if err != nil {
			return
		}
	}
	_, err = w.WriteString("?>")
	return
}

func (e *XMLDirective) Encode(w ByteAndStringWriter) (err error) {
	w.WriteString("<!")
	w.Write(e.Directive)
	w.WriteString(">")
	return
}
