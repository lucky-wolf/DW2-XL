package main

import (
	"flag"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
)

var (
	source string
	output string
	data   string
)

func main() {

	flag.StringVar(&source, "target", "", "specifies the file to apply the changes to")
	// flag.StringVar(&data, "data", "", "specifies the data file to extract new values from")
	flag.Parse()
	output = source + ".bak"

	log.SetFlags(0)

	path, err := os.Getwd()
	if err != nil {
		return
	}

	log.Printf("cwd=%s", path)

	err = run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {

	// // build the new map of race id -> factors
	// actions, err := load()
	// if err != nil {
	// 	return
	// }

	// process our input + data -> output
	tree, err := xmltree.LoadFromFile(source)
	if err != nil {
		return
	}

	// debug: just print out the tree
	// log.Print(tree)

	// bytes, err := xml.MarshalIndent(tree, "", "\t")
	// if err != nil {
	// 	return
	// }

	// log.Print(string(bytes))

	e := xmltree.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	err = tree.Encode(e)

	// // finally, swap the new races files so it's ready to go
	// err = swap(oldfile, newfile)
	// if err != nil {
	// 	return
	// }

	return
}
