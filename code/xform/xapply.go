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
	for i := range xs.Values {
		stats.Searched++
		_, e := tree.Find(searchTag, xs.Values[i][searchColumn])
		if e == nil {
			log.Printf("warn: did not find %s = %s", searchTag, xs.Values[i][searchColumn])
			continue
		}
		stats.Found++
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
