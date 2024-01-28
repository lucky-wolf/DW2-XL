package xmltree

import (
	"errors"
	"fmt"
	"regexp"
)

// clone an element
func (e *XMLElement) Clone() *XMLElement {
	return &XMLElement{
		StartElement: e.StartElement.Copy(),
		XMLValue:     e.XMLValue.Clone(),
	}
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

// returns the first matching element from the list of elements based on regex of tag-name
func (e *XMLElement) Matching(r *regexp.Regexp) (children []*XMLElement) {

	// scan our elements for those that match
	for _, e = range e.Elements() {
		if r.MatchString(e.Name.Local) {
			children = append(children, e)
		}
	}

	return
}

// returns the first matching child element whose tag and value equal the find tag and value
func (e *XMLElement) FindRecurse(tag string, value string) *XMLElement {
	for _, e = range e.Elements() {
		if e.Name.Local == tag && e.StringValueEquals(value) {
			return e
		}
		f := e.FindRecurse(tag, value)
		if f != nil {
			return f
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

var ErrTargetNodeNotFound = errors.New("target doesn't have a node to copy to")

// copies the specified child from the given element to the ourself (replacing any we already have)
func (e *XMLElement) CopyByTag(tag string, from *XMLElement) (err error) {

	// get source
	source := from.Child(tag)
	if source == nil {
		err = fmt.Errorf("%s doesn't have a %s to copy from", from.Child("Name").StringValue(), tag)
		return
	}

	// get target
	target := e.Child(tag)
	if target == nil {
		err = fmt.Errorf("%s doesn't have a %s to copy from", e.Child("Name").StringValue(), tag)
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
	err = e.CopyByTag(tag, from)
	if err != nil {
		return
	}

	// visit the new elements
	err = e.Child(tag).VisitChildren(visit)
	if err != nil {
		return
	}

	return
}

// updates it to be scaled by the given input
func (e *XMLElement) ScaleChildBy(tag string, scale float64) (err error) {

	if scale == 1.0 {
		return
	}

	return e.Child(tag).ScaleBy(scale)
}

// sets one value to be that of another (both must be value types)
func (e *XMLElement) SetChildToSibling(child, sibling string) {
	e.Child(child).SetValue(e.Child(sibling).StringValue())
}

// updates it to be scaled by the given input
func (e *XMLElement) ScaleChildToSiblingBy(child, sibling string, scale float64) {
	value := e.Child(sibling).NumericValue()
	e.Child(child).SetValue(value * scale)
}

// updates it to be scaled by the given input
func (e *XMLElement) AdjustChildToSiblingBy(child, sibling string, adj float64) {
	value := e.Child(sibling).NumericValue()
	e.Child(child).SetValue(value + adj)
}
