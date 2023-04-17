package main

// this works if we didn't care about order
// but, because we want to compare between files, in general
// we do care - and this sorts each set of tags rather than keeping the original ordering

import (
	"flag"
	"log"
	"os"

	"github.com/clbanning/mxj/v2"
)

var (
	source string
	target string
)

func main() {

	log.SetFlags(0)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(cwd)

	flag.StringVar(&source, "source", "", "source filename")
	flag.StringVar(&target, "target", "", "target filename")
	flag.Parse()

	err = run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {

	// any xml file should be readable
	values, err := ReadXMLFile(source)
	if err != nil {
		return
	}

	// just dump the file to stdout for now
	// fmt.Println(values.StringIndentNoTypeInfo())

	// convert it to output
	err = WriteXMLFile(target, values)

	return
}

func ReadXMLFile(filename string) (m mxj.Map, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	mxj.CastValuesToInt()
	mxj.CastValuesToBool()
	m, err = mxj.NewMapXmlReader(file) // file, true -- uses the above casts
	if err != nil {
		return
	}

	return
}

func WriteXMLFile(filename string, m mxj.Map) (err error) {

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	err = m.XmlIndentWriter(file, "", "\t")
	// err = m.XmlWriter(file)

	return
}
