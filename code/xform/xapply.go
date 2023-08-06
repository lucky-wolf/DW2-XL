package xform

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

type Stats struct {
	Searched int
	Found    int
	Replaced int
}

func (stats Stats) String() string {
	return fmt.Sprintf("Searched: %d, Found: %d, Replaced: %d", stats.Searched, stats.Found, stats.Replaced)
}

func (xs XScript) ApplyTo(tree *xmltree.XMLTree) (stats Stats, err error) {

	// basically we want to scan the tree (iterator?) for each search term...
	// then apply the changes

	searchTag := xs.Search[0]
	searchColumn := xs.indexOf(searchTag)
	searchTag = searchTag[1 : len(searchTag)-1]

	// this is a very specific way of encoding the transformations, so we'll keep the logic here...
	var targetValue string
	var e *xmltree.XMLElement
	var subindex int
	for r := 0; r < len(xs.Values); r++ {

		stats.Searched++

		// try to find it
		if len(xs.Values[r][searchColumn]) == 0 {
			// if the search value is blank, we assume we're modifying the current object, on a further subelement
			subindex++
		} else {
			// non-blank search value: search for it explicitly
			targetValue = xs.Values[r][searchColumn]
			_, e = tree.Find(searchTag, targetValue)
			if e == nil {
				log.Printf("warn: did not find %s = %s", searchTag, targetValue)
				continue
			}
		}
		stats.Found++

		// now to build a way to apply the transform data to the node we found (and its underlying tiers)
		// for each column in our values, apply it
		// for c := range xs.Header {

		// }

		//todo: we'll need to apply rows with no id as additional component tiers.
	}

	return
}

func (xs XScript) indexOf(tag string) int {
	for i := range xs.Header {
		if xs.Header[i] == tag {
			return i
		}
	}
	return -1
}
