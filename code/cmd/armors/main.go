package main

import (
	"flag"
	"log"
	"lucky-wolf/DW2-XL/code/xform"
	"lucky-wolf/DW2-XL/code/xmltree"
	"os"
	"time"
)

var (
	source string
	target string
	script string
)

func main() {

	flag.StringVar(&source, "source", "", "input filename")
	flag.StringVar(&target, "target", "", "target filename")
	flag.StringVar(&script, "script", "", "script filename")
	flag.Parse()

	log.SetFlags(0)

	path, err := os.Getwd()
	if err != nil {
		return
	}

	log.Printf("cwd=%s", path)

	t := time.Now()
	err = run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("time=%s", time.Since(t).String())
}

func run() (err error) {

	// load our transform script
	actions, err := xform.LoadFromFile(script)
	if err != nil {
		return
	}

	log.Print(actions)

	// process our input + data -> target
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

	err = tree.WriteToFile(target)
	// e := xmltree.NewEncoder(os.Stdout)
	// e.Configure("", "\t")
	// err = tree.Encode(e)

	// // finally, swap the new races files so it's ready to go
	// err = swap(oldfile, newfile)
	// if err != nil {
	// 	return
	// }

	return
}
