package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
)

type Transformation = func(data *xmltree.XMLTree) (statistics Statistics, err error)

type Statistics struct {
	objects  uint
	elements uint
	changed  uint
}

var (
	Quiet bool
)

func (statistics *Statistics) For(filename string) string {
	return fmt.Sprintf("%s: objects found: %d, elements updated: %d of %d", filename, statistics.objects, statistics.changed, statistics.elements)
}

func AssertIs(e *xmltree.XMLElement, kind string) (err error) {
	if !Quiet {
		log.Println(e.Name.Local)
	}
	if e.Name.Local != kind {
		err = fmt.Errorf("invalid file: expected %s but found %s", kind, e.Name.Local)
	}
	return
}
