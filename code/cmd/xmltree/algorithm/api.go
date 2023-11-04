package algorithm

import (
	"fmt"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
	"path/filepath"
	"strings"
)

var (
	Quiet bool
)

type Statistics struct {
	objects  int
	elements int
	changed  int
}

type LevelFunc = func(level int) float64
type ValuesTable = map[string]LevelFunc

func ExtendValuesTable(fields, more ValuesTable) (result ValuesTable) {
	result = ValuesTable{}
	for k, v := range fields {
		result[k] = v
	}
	for k, v := range more {
		result[k] = v
	}
	return
}

type ComponentData struct {
	// todo: could we ever look this up viz research tree for first occurrence there?
	//       that's just column in which it is listed for each level...
	scaleToFtr  bool
	minLevel    int
	maxLevel    int
	fieldValues ValuesTable
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

type Job struct {
	xfiles []*XFile
}

type XFile struct {
	filename string
	root     *xmltree.XMLTree
	stats    Statistics
}

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

func LoadJobFor(root, pattern string) (j Job, err error) {

	// get the list of files applicable
	filenames, err := FindMatchingFiles(root, pattern)
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

// apply stats for each component
func (j *Job) ApplyAll(components map[string]ComponentData) (err error) {
	for k, v := range components {
		err = j.Apply(k, v)
		if err != nil {
			return
		}
	}
	return
}

// applies stats for given component
func (j *Job) Apply(name string, data ComponentData) (err error) {

	// find this component definition
	e, f := j.FindElement("Name", name)
	if e == nil {
		return fmt.Errorf("%s not found", name)
	}
	statistics := &f.stats

	// ensure we have correct number of component stats to update
	err = e.Child("Values").SetElementCount(1 + data.maxLevel - data.minLevel)
	if err != nil {
		return
	}

	// fill in the data from our data tables
	stats := e.Child("Values").Elements()
	for i, e := range stats {
		for key, f := range data.fieldValues {
			e.Child(key).SetValue(f(data.minLevel + i))
		}
		statistics.elements++
		statistics.changed++
	}
	statistics.objects++

	// scale to fighter if required
	if data.scaleToFtr {
		err = j.ScaleToFighter(e)
	}

	return
}

func (j *Job) ScaleToFighter(sourceDefinition *xmltree.XMLElement) (err error) {

	name := sourceDefinition.Child("Name").StringValue()
	if strings.HasSuffix(name, " [S]") || strings.HasSuffix(name, " [M]") || strings.HasSuffix(name, " [L]") {
		name = name[:len(name)-3]
	}

	// apply to [ftr]
	name += " [Ftr]"
	e, f := j.FindElement("Name", name)
	if e == nil {
		err = fmt.Errorf("%s not found", name)
		return
	}
	statistics := &f.stats

	// copy and scale resource requirements
	err = e.CopyAndVisitByTag("ResourcesRequired", sourceDefinition, func(e *xmltree.XMLElement) error { e.Child("Amount").ScaleBy(0.2); return nil })
	if err != nil {
		log.Println(err)
	}

	// copy component stats
	err = e.CopyByTag("Values", sourceDefinition)
	if err != nil {
		log.Println(err)
	}

	// now that we have our own copy of the component stats (same number of levels too)
	// we can update each of those to scale for [Ftr] version
	for _, e := range e.Child("Values").Elements() {

		// every element should be a component bay
		err = assertIs(e, "ComponentStats")
		if err != nil {
			return
		}

		// scale down the ion defenses (or ion PD / ftr weapons will never penetrate)
		// e.Child("ComponentIonDefense").ScaleBy(1)
		// e.Child("IonDamageDefense").ScaleBy(1)

		// never a crew requirement for fighter components
		e.Child("CrewRequirement").SetValue(0)

		// 25% static draw
		e.Child("StaticEnergyUsed").SetValue(0.25)

		statistics.changed++
		statistics.elements++
	}

	statistics.objects++

	return
}
