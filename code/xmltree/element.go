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

// like Find, but looks for anything with the value that matches as a prefix
func (e *XMLElement) HasPrefix(tag string, prefix string) bool {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueStartsWith(prefix) {
			return true
		}
	}
	return false
}

// like Find, but looks for anything with the value that matches as a suffix
func (e *XMLElement) HasSuffix(tag string, suffix string) bool {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueEndsWith(suffix) {
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

	case *XMLElement:
		e := v.Clone()
		return &e
	case XMLComment:
		return v.Copy()
	case XMLDirective:
		return v.Copy()
	case XMLElement:
		return v.Clone()
	case XMLProcInst:
		return v.Copy()
	case string:
		return v
	}
	panic(fmt.Errorf("cannot clone: invalid contents: %T", contents))
}

// visits each child with the given visitor function (aborts on error)
func (e *XMLElement) VisitChildren(visit func(*XMLElement) error) (err error) {
	if e == nil {
		return
	}
	for _, e := range e.Elements() {
		err = visit(e)
		if err != nil {
			return
		}
	}
	return
}

// copies the specified child from the given element to the ourself (replacing any we already have)
func (e *XMLElement) CopyByTag(tag string, from *XMLElement) (source, target *XMLElement, err error) {

	// get source
	source = from.Child(tag)
	if source == nil {
		// err = fmt.Errorf("source doesn't have a <%s> node!", tag)
		return
	}

	// get target
	target = e.Child(tag)
	if target == nil {
		err = fmt.Errorf("target doesn't have a <%s> node!", tag)
		return
	}

	// clone source to target
	target.SetContents(source.CloneContents())

	return
}

// copies the specified child from the given element to the ourself (replacing any we already have)
// and visits the new elements
func (e *XMLElement) CopyAndVisitByTag(tag string, from *XMLElement, visit func(*XMLElement) error) (err error) {

	// copy by tag
	_, target, err := e.CopyByTag(tag, from)
	if err != nil {
		return
	}

	// visit the new elements
	target.VisitChildren(visit)

	return
}
