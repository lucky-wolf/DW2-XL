package xmltree

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/cmd/etc"
)

// returns XMLElements only
func (v *XMLValue) Elements() (elements []*XMLElement) {

	// it is legal to call on a nil value (we simply have no child elements)
	if v == nil {
		return
	}

	// we are either a single or multiple elements
	switch v := v.contents.(type) {
	case *XMLElement:
		// single element
		elements = append(elements, v)
	case []any:
		// multiple child elements
		for _, e := range v {
			switch v := e.(type) {
			case *XMLElement:
				elements = append(elements, v)
			}
		}
	default:
		// we have no child elements
	}

	return
}

func (v XMLValue) Clone() XMLValue {
	// v is already a shallow copy, just do a deep copy on the contents
	v.contents = CloneContents(v.contents)
	return v
}

func (v *XMLValue) SetContents(contents any) {
	v.contents = contents
	return
}

func (v *XMLValue) CloneContents() any {
	return CloneContents(v.contents)
}

func CloneContents(contents any) any {

	switch t := contents.(type) {

	case []any:
		// multiple child contents
		contents := []any{}
		for _, e := range t {
			contents = append(contents, CloneContents(e))
		}
		return contents

	case *XMLElement:
		return &XMLElement{StartElement: t.StartElement.Copy(), XMLValue: t.XMLValue.Clone()}
	case *XMLComment:
		return &XMLComment{Comment: t.Copy()}
	case *XMLDirective:
		return &XMLDirective{Directive: t.Copy()}
	case *XMLProcInst:
		return &XMLProcInst{ProcInst: t.Copy()}

	// case XMLElement:
	// 	return XMLElement{StartElement: v.StartElement.Copy(), XMLValue: v.XMLValue.Clone()}
	// case XMLComment:
	// 	return XMLComment{Comment: v.Copy()}
	// case XMLDirective:
	// 	return XMLDirective{Directive: v.Copy()}
	// case XMLProcInst:
	// 	return XMLProcInst{ProcInst: v.Copy()}

	case string:
		return t
	}

	err := fmt.Errorf("cannot clone: invalid contents: %T", contents)
	log.Fatal(err)
	panic(err)
}

// ensures we have count copies of our first element
func (v *XMLValue) SetElementCountByCopyingFirstElementAsNeeded(count int) (err error) {

	elements := v.Elements()
	l := len(elements)
	switch {

	case l == 0:
		return fmt.Errorf("no elements found")

	case l < count:
		log.Printf("extending by %d elements", count-l)
		for i := l; i < count; i++ {
			v.Append(elements[0])
		}

	case l > count:
		log.Printf("truncating by %d elements", l-count)
		v.Truncate(count)
	}

	return
}

// we must already be a []any or we error
func (v *XMLValue) Append(e *XMLElement) (err error) {

	switch t := v.contents.(type) {

	case []any:
		v.contents = append(t, e.Clone())

	default:
		// v.contents = []any{t, e.Clone()}
		err = fmt.Errorf("append requires []any")
	}

	return
}

// we must already be a []any or this is an error
func (v *XMLValue) Truncate(count int) (err error) {

	switch t := v.contents.(type) {

	case []any:
		v.contents = t[:count]

	default:
		// v.contents = []any{t, e.Clone()}
		err = fmt.Errorf("truncate requires []any")
	}

	return
}

// we must already be a []any or we error
func (v *XMLValue) InsertCopyOf(index, copy int) (err error) {

	switch t := v.contents.(type) {

	case []any:
		t = etc.InsertAt(t, index, t[copy])

	default:
		// v.contents = []any{t, e.Clone()}
		err = fmt.Errorf("append requires []any")
	}

	return
}

// we must already be a []any or we error
// WARN! we'll use the element you hand us, if you need to copy it, use e.Clone() when calling us!
func (v *XMLValue) InsertAt(index int, e *XMLElement) (err error) {

	switch t := v.contents.(type) {

	case []any:
		t = etc.InsertAt(t, index, any(e))

	default:
		// v.contents = []any{t, e.Clone()}
		err = fmt.Errorf("append requires []any")
	}

	return
}
