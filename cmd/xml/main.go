package main

import (
	"fmt"
	"log"
	"os"

	"github.com/clbanning/mxj/v2"
)

func main() {

	log.SetFlags(0)

	log.Println(os.Getwd())

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {

	// C:\\Users\\steve\\Projects\\DW2-XL
	values, err := ReadXMLFile("../XL/ShipHulls_Ackdarian.xml")
	if err != nil {
		return
	}

	// just dump the file to stdout for now
	// fmt.Println(values.StringIndent())

	fmt.Println(values.Xml())

	return
}

func ReadXMLFile(filename string) (m mxj.Map, err error) {

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	mxj.CastValuesToInt()
	m, _, err = mxj.NewMapXmlReaderRaw(file, true)
	if err != nil {
		return
	}

	return
}
