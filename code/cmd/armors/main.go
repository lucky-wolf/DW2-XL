package main

import (
	"flag"
	"log"
	"lucky-wolf/DW2-XL/code/xform"
	"lucky-wolf/DW2-XL/code/xmltree"
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

	// log our cwd
	// path, err := os.Getwd()
	// if err != nil {
	// 	return
	// }
	// log.Printf("cwd=%s", path)

	t := time.Now()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("time=%s", time.Since(t).String())
}

func run() (err error) {

	log.Print("source: ", source)
	log.Print("target: ", target)
	log.Print("script: ", script)

	// load our transform script
	actions, err := xform.LoadFromFile(script)
	if err != nil {
		return
	}

	// debug: log the action script
	// log.Print(actions)

	// process our input + data -> target
	tree, err := xmltree.LoadFromFile(source)
	if err != nil {
		return
	}

	// attempt to apply our actions to the data tree
	stats, err := actions.ApplyTo(tree)
	if err != nil {
		return
	}

	// log the search & replace stats
	log.Print(stats)

	// debug: log the xml tree
	// log.Print(tree)

	// write out the new file
	err = tree.WriteToFile(target)

	// // finally, swap the new races files so it's ready to go
	// err = swap(oldfile, newfile)
	// if err != nil {
	// 	return
	// }

	return
}
