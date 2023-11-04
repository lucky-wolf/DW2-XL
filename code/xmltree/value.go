package xmltree

import (
	"fmt"
	"log"
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

// we are converted to an []any if we aren't already
// note: we append a clone of the specified element
func (v *XMLValue) Append(e *XMLElement) {

	switch t := v.contents.(type) {

	case string:
		err := fmt.Errorf("cannot append to a string XMLValue")
		log.Fatal(err)
		panic(err)

	case []any:
		v.contents = append(t, e.Clone())

	default:
		v.contents = []any{t, e.Clone()}
	}
}
