package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
	"path/filepath"
)

// global to control log output verbosity
var (
	Quiet bool
)

// The Job is the base object for building an algorithm
// it holds the list of files to be updated and has a bunch of helper members for manipulating them
type Job struct {
	xfiles []*XFile
}

// One loaded xml file + update statistics for that file
type XFile struct {
	filename string
	root     *xmltree.XMLTree
	stats    Statistics
}

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

/////////////////////////////////////////////////////////////////////////

func FindMatchingFiles(root, pattern string) (matches []string, err error) {

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

func LoadJobFor(root string, patterns ...string) (j Job, err error) {

	for _, pattern := range patterns {

		// get the list of files applicable
		var filenames []string
		filenames, err = FindMatchingFiles(root, pattern)
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
			j.xfiles = append(j.xfiles, &XFile{filename: filenames[i], root: root})
		}
	}

	return
}

func (j *Job) Save() {
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

func (j *Job) FindElement(tag, value string) (element *xmltree.XMLElement, file *XFile) {
	for _, file = range j.xfiles {
		element, _ = file.root.Find(tag, value)
		if element != nil {
			return
		}
	}
	return
}
