package xmltree

import (
	"encoding/xml"
	"io"
	"os"
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

	// use default formating for a file
	encoder.Configure("", "\t")

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
