package main

// trivial test harness for our custom xmltree lib
// simply want to read & write without scrambling the file or losing comments

import (
	"flag"
	"log"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
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
	values, err := xmltree.LoadFromFile(source)
	if err != nil {
		return
	}

	// convert it to output
	err = values.WriteToFile(target)

	return
}
