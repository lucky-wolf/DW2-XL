package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
	"path/filepath"
)

var (
	Quiet bool
)

type Statistics struct {
	objects  int
	elements int
	changed  int
}

func (statistics *Statistics) For(filename string) string {
	return fmt.Sprintf("%s: objects found: %d, elements updated: %d of %d", filename, statistics.objects, statistics.changed, statistics.elements)
}

func assertIs(e *xmltree.XMLElement, kind string) (err error) {
	// if !Quiet {
	// 	log.Println(e.Name.Local)
	// }
	if e.Name.Local != kind {
		err = fmt.Errorf("invalid file: expected %s but found %s", kind, e.Name.Local)
	}
	return
}

type job struct {
	xfiles []*xmlfile
}

type xmlfile struct {
	filename string
	root     *xmltree.XMLTree
	stats    Statistics
}

func findMatchingFiles(root, pattern string) (matches []string, err error) {

	walker := func(path string, fi os.FileInfo, pe error) (err error) {

		if pe != nil || fi.IsDir() {
			return
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return
		}

		if matched {
			matches = append(matches, path)
		}

		return
	}

	err = filepath.Walk(root, walker)
	return
}

func loadJobFor(root, pattern string) (j job, err error) {

	// get the list of files applicable
	filenames, err := findMatchingFiles(root, pattern)
	if err != nil {
		return
	}

	// load each and every one so we have access to all of them
	for i := range filenames {
		var root *xmltree.XMLTree
		root, err = xmltree.LoadFromFile(filenames[i])
		if err != nil {
			return
		}
		j.xfiles = append(j.xfiles, &xmlfile{filename: filenames[i], root: root})
	}

	return
}

func (j *job) save() {
	// save all of our files
	for _, f := range j.xfiles {
		err := f.root.WriteToFile(f.filename)
		switch err {
		case nil:
			log.Println(f.stats.For(f.filename))
		default:
			log.Printf("failed to write %s: %s", f.filename, err)
		}
	}
}

func (j *job) find(tag, value string) (parent, element *xmltree.XMLElement) {
	for _, f := range j.xfiles {
		parent, element = f.root.Find(tag, value)
		if parent != nil {
			return
		}
	}
	return
}
