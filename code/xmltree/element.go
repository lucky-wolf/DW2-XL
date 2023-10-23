package xmltree

import "fmt"

// element member functions

// returns XMLElements only
func (e *XMLValue) Elements() (elements []*XMLElement) {

	// it is legal to call on a nil value (we simply have no child elements)
	if e == nil {
		return
	}

	// we are either a single or multiple elements
	switch v := e.contents.(type) {
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

// returns the first matching element from the list of elements based on tag (name)
func (e *XMLElement) Child(tag string) *XMLElement {
	for _, e = range e.Elements() {
		if e.Name.Local == tag {
			return e
		}
	}
	return nil
}

// returns the first matching child element whose tag and value equal the find tag and value
func (e *XMLElement) Find(tag string, value string) *XMLElement {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueEquals(value) {
			return e
		}
	}
	return nil
}

// returns true if the given element has a sub element with specified tag and value
func (e *XMLElement) Has(tag string, value string) bool {
	return e.Find(tag, value) != nil
}

// like Find, but looks for anything with the value that matches as a suffix
func (e *XMLElement) HasEndsWith(tag string, tail string) bool {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueEquals(tail) {
			return true
		}
	}
	return false
}

// clones the given element (always a deep copy)
func (e *XMLElement) Clone() (c XMLElement) {

	// nil -> nil
	if e == nil {
		return
	}

	// copy our start element
	c.StartElement = e.StartElement.Copy()

	// copy our contents
	c.XMLValue = e.XMLValue.Clone()

	return
}

func (e *XMLValue) Clone() (v XMLValue) {
	v.contents = CloneContents(e.contents)
	return
}

func (e *XMLValue) SetContents(contents any) {
	e.contents = contents
	return
}

func (e *XMLValue) CloneContents() any {
	return CloneContents(e.contents)
}

func CloneContents(contents any) any {
	switch v := contents.(type) {
	case []any:
		// multiple child elements
		elements := []any{}
		for _, e := range v {
			elements = append(elements, CloneContents(e))
		}
		return elements
	case XMLProcInst:
		return v.Copy()
	case XMLComment:
		return v.Copy()
	case XMLDirective:
		return v.Copy()
	case XMLElement:
		return v.Clone()
	case string:
		return v
	}
	panic(fmt.Errorf("cannot clone: invalid contents: %T", contents))
}
