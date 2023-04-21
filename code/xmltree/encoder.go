package xmltree

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// we encode a tree to a stream or file

var ErrClosed = errors.New("use of closed Encoder")

func UnknownEntity(e any) error {
	return fmt.Errorf("unknown entity: %T", e)
}

// required capabilities
type ByteAndStringWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
}
type FormattedEncoder interface {
	ByteAndStringWriter
	io.Closer
	Indent(newline bool, changeDepth int, indent bool) (err error)
	Flush() error
}

// represents our config choices for outputting an xml tree to a stream
type encoder struct {
	writer *bufio.Writer
	prefix string
	indent string
	depth  int
	closed bool
}

// NewEncoder returns a new encoder that writes to w
func NewEncoder(stream io.Writer) (e *encoder) {
	e = &encoder{writer: bufio.NewWriter(stream)}
	return
}

// Configures this encoder to use the given prefix + indent
func (e *encoder) Configure(prefix, indent string) {
	e.prefix = prefix
	e.indent = indent
	return
}

// Flushes any buffered XML to the underlying writer
func (e *encoder) Flush() (err error) {
	err = e.writer.Flush()
	return
}

// Closes the encoder, indicating that no more data will be written
// It flushes any buffered XML to the underlying writer
// implements io.Closer
func (e *encoder) Close() (err error) {
	if e.closed {
		return
	}
	e.closed = true
	err = e.writer.Flush()
	return
}

// 1. optionally terminates the current line
// 2. updates current depth
// 3. writes our indent string x new depth
func (p *encoder) Indent(newline bool, changeDepth int, indent bool) (err error) {

	// terminate the current line, and write prefix + indent for start of new line
	if newline {
		err = p.WriteByte('\n')
		if err != nil {
			return
		}
	}

	// sanity check the new depth before applying it
	newdepth := p.depth + changeDepth
	if newdepth < 0 {
		err = fmt.Errorf("negative indentation depth")
		return
	}

	// update our depth
	p.depth = newdepth

	if indent {
		// always write the prefix
		if len(p.prefix) > 0 {
			_, err = p.WriteString(p.prefix)
			if err != nil {
				return
			}
		}

		// write as many indents as our current depth requires
		if len(p.indent) > 0 {
			for i := 0; i < p.depth; i++ {
				_, err = p.WriteString(p.indent)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// Write implements io.Writer
func (e *encoder) Write(b []byte) (n int, err error) {
	if e.closed {
		err = ErrClosed
		return
	}
	n, err = e.writer.Write(b)
	return
}

// WriteString implements io.StringWriter
func (e *encoder) WriteString(s string) (n int, err error) {
	if e.closed {
		err = ErrClosed
		return
	}
	n, err = e.writer.WriteString(s)
	return
}

// WriteByte implements io.ByteWriter
func (e *encoder) WriteByte(c byte) (err error) {
	if e.closed {
		err = ErrClosed
		return
	}
	err = e.writer.WriteByte(c)
	return
}

func QuotedString(value string) string {
	return `"` + EscapeString(value) + `"`
}

var (
	escQuote = []byte("&quot;")
	escTick  = []byte("&apos;")
	escAmp   = []byte("&amp;")
	escLT    = []byte("&lt;")
	escGT    = []byte("&gt;")
	escTab   = []byte("&#x9;")
	escNL    = []byte("&#xA;")
	escCR    = []byte("&#xD;")
	escFF    = []byte("\uFFFD") // Unicode replacement character
)

// full set to escape within a string
var strictMapping = map[rune][]byte{
	'"':  escQuote,
	'\'': escTick,
	'&':  escAmp,
	'<':  escLT,
	'>':  escGT,
	'\t': escTab,
	'\n': escNL,
	'\r': escCR,
}

// only < and & are strictly illegal in XML
var looseMapping = map[rune][]byte{
	'&': escAmp,
	'<': escLT,
}

// EscapeString returns the properly escaped XML equivalent of the plain text data s
func EscapeString(s string) string {
	sb := &strings.Builder{}
	WriteEscapedText(s, sb, true)
	return sb.String()
}

func WriteEscapedText(s string, sb ByteAndStringWriter, strict bool) (err error) {

	// choose the strict or loose character mapping
	var mapping map[rune][]byte
	if strict {
		mapping = strictMapping
	} else {
		mapping = looseMapping
	}

	// walk the input string substituting as we go
	last := 0
	var esc []byte
	for i := 0; i < len(s); {
		// get next rune
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width

		// check if we have an esc mapping for this rune
		esc = mapping[r]
		if len(esc) == 0 && (!IsInCharacterRange(r) || (r == 0xFFFD && width == 1)) {
			esc = escFF
		} else if len(esc) == 0 {
			// no mapping: continue
			continue
		}

		// we have an esc substitution, so write up to but not including current rune
		_, err = sb.WriteString(s[last : i-width])
		if err != nil {
			return
		}

		// and write the substitution
		_, err = sb.Write(esc)
		if err != nil {
			return
		}

		// start the next span
		last = i
	}

	// write any trailing bit
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
