package xmltree

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
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

// true if we hold nothing (we're the empty value)
func (e *XMLValue) Empty() bool {

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

// true if we hold a simple string value
func (e *XMLValue) IsSimple() bool {
	_, ok := e.contents.(string)
	return ok
}

// true if we hold multiple children
func (e *XMLValue) HasMultipleChildren() bool {
	_, ok := e.contents.([]any)
	return ok
}

// returns the string value of this value iff it is a simple value
func (e *XMLValue) GetStringValue() (s string, ok bool) {
	s, ok = e.contents.(string)
	return
}

// returns the string value of this value
// panics if it is not a string
func (e *XMLValue) StringValue() string {
	return e.contents.(string)
}

// returns the string value of this value iff it is a simple value
func (e *XMLValue) StringValueEquals(value string) bool {
	s, ok := e.contents.(string)
	return ok && s == value
}

// returns the string value of this value iff it is a simple value
func (e *XMLValue) StringValueStartsWith(value string) bool {
	s, ok := e.contents.(string)
	return ok && strings.HasPrefix(s, value)
}

// returns the string value of this value iff it is a simple value
func (e *XMLValue) StringValueEndsWith(value string) bool {
	s, ok := e.contents.(string)
	return ok && strings.HasSuffix(s, value)
}

// set our contents to the given string value-string
func (e *XMLValue) SetString(value string) {
	_, ok := e.contents.(string)
	if !ok {
		panic("not a simple value type: cannot write a simple value into it")
	}
	e.contents = value
}

// set our contents to the given value
// value can be any kind of scalar or string
func (e *XMLValue) SetValue(value any) {
	switch v := value.(type) {
	case string:
		e.SetString(v)
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8:
		e.SetString(fmt.Sprint(v))
	case float32, float64:
		e.SetString(fmt.Sprintf("%.6g", v))
	default:
		e.SetString(fmt.Sprint(v))
	}
}

// if the value is simple and parsable as float, returns that
func (e *XMLValue) GetNumericValue() (value float64, err error) {
	// must be simple
	s, ok := e.contents.(string)
	if !ok {
		err = fmt.Errorf("XMLValue is not simple: cannot extract a value from it")
		return
	}

	// must be parsable as a float
	return strconv.ParseFloat(s, 64)
}

// grab string & parse (may end up hiding errors and being zero)
func (e *XMLValue) NumericValue() (value float64) {
	value, _ = strconv.ParseFloat(e.StringValue(), 64)
	return
}

// our contents must be a simple string which is a parsable number
// updates it to be scaled by the given input
func (e *XMLValue) ScaleBy(scale float64) (err error) {

	value, err := e.GetNumericValue()
	if err != nil {
		return
	}

	e.contents = fmt.Sprintf("%.5g", value*scale)
	return
}

// updates it to be current value + adjustment
func (e *XMLValue) AdjustValue(adjustment float64) (err error) {

	value, err := e.GetNumericValue()
	if err != nil {
		return
	}

	e.contents = fmt.Sprintf("%.5g", value+adjustment)
	return
}

////////////////////////////////////////////////////
// simple string representation

func (e *XMLTree) String() string {
	return e.Elements.String()
}

func (e *XMLValue) String() string {

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
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (e *XMLElement) String() string {

	sb := new(strings.Builder)

	// write the start token with attributes
	sb.WriteByte('<')
	sb.WriteString(e.Name.Local)
	for i := range e.Attr {
		sb.WriteByte(' ')
		if e.Attr[i].Name.Space != "" {
			sb.WriteString(e.Attr[i].Name.Space)
			sb.WriteByte(':')
		}
		sb.WriteString(e.Attr[i].Name.Local)
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

func (e *XMLProcInst) String() string {
	sb := &strings.Builder{}
	e.Encode(sb)
	return sb.String()
}

func (e *XMLDirective) String() string {
	sb := &strings.Builder{}
	e.Encode(sb)
	return sb.String()
}

func (e *XMLComment) String() string {
	sb := &strings.Builder{}
	e.Encode(sb)
	return sb.String()
}
