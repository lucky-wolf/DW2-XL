package xmltree

type ParentChildElement struct {
	Parent *XMLElement
	Child  *XMLElement
}

// find the given element which has that value
// this algo only finds the first one (if there are multiple)
// returns the parent and the element (parent is nil if th root element is the matching target)
func (tree *XMLTree) Find(tag, value string) (parent, element *XMLElement) {

	// use a simple bfs algo
	var queue []*ParentChildElement
	pop := func() *ParentChildElement {
		if len(queue) == 0 {
			return nil
		}
		e := queue[0]
		queue = queue[1:]
		return e
	}

	// start with our (root) element(s) (which have no parent)
	for _, e := range tree.Elements.Elements() {
		queue = append(queue, &ParentChildElement{nil, e})
	}

	// process each node until we find it, or there are no more nodes
	for pc := pop(); pc != nil; pc = pop() {

		// check if we've found our target
		if pc.Child.Name.Local == tag {
			v, ok := pc.Child.StringValue()
			if ok && v == value {
				parent = pc.Parent
				element = pc.Child
				return
			}
		}

		// queue up this child's children
		for _, e := range pc.Child.Elements() {
			queue = append(queue, &ParentChildElement{pc.Child, e})
		}
	}

	return
}

// returns XMLElements only
func (e *XMLValue) Elements() (elements []*XMLElement) {
	switch v := e.contents.(type) {
	case *XMLElement:
		elements = append(elements, v)
	case []any:
		for _, e := range v {
			switch v := e.(type) {
			case *XMLElement:
				elements = append(elements, v)
			}
		}
	}
	return
}
