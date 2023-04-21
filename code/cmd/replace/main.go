package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// go run cmd/replace/main.go -script renamehulls.txt -target XL/ShipHulls_Ackdarian.xml
// target will be backed up as .bak
// replace must contain a search pattern, followed by a replacement pattern, whitespace separated

var (
	script string
	source string
	output string
)

func main() {

	var err error

	log.SetFlags(0)

	flag.StringVar(&source, "target", "", "specifies the file to apply the changes to")
	flag.StringVar(&script, "script", "", "specifies the file to extract search and replacement strings from")
	flag.Parse()

	output = source + ".bak"

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

	// build the new map of race id -> factors
	actions, err := load()
	if err != nil {
		return
	}

	// process our input + data -> output
	err = execute(actions)
	if err != nil {
		return
	}

	// finally, swap the new races files so it's ready to go
	err = swap(source, output)
	if err != nil {
		return
	}

	return
}

// loads our script file
func load() (actions map[string]string, err error) {

	file, err := os.Open(script)
	if err != nil {
		return
	}
	defer file.Close()

	// handle the header line
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		err = fmt.Errorf("empty file")
		return
	}

	// process this line by line
	actions = map[string]string{}
	for line := 1; scanner.Scan(); line++ {

		text := strings.TrimSpace(scanner.Text())

		// skip blanks
		if text == "" {
			continue
		}

		// skip comments
		if text[:2] == "//" {
			continue
		}

		// each line should contain the search and replacement strings
		values := split(text)
		if len(values) != 2 {
			err = fmt.Errorf("invalid line %d: does not contain 2 fields", line)
			return
		}

		// record them
		actions[values[0]] = values[1]
	}

	return
}

func split(text string) (values []string) {
	i := strings.Index(text, "\t")
	if i == -1 {
		return
	}
	j := i
	for text[j] == '\t' {
		j++
		if j == len(text) {
			return
		}
	}
	values = append(values, text[:i])
	values = append(values, text[j:])
	return
}

func execute(actions map[string]string) (err error) {

	// open our races.xml as our input
	input, err := os.Open(source)
	if err != nil {
		return
	}
	defer input.Close()

	output, err := os.Create(output)
	if err != nil {
		return
	}
	defer output.Close()

	// create a buffered writer
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// scan line by line
	scanner := bufio.NewScanner(input)

	for lineIndex := 1; scanner.Scan(); lineIndex++ {

		// grab a line to process
		line := scanner.Text()

		for s, r := range actions {
			i := strings.Index(line, s)
			if i != -1 {
				line = strings.ReplaceAll(line, s, r)
				break
			}
		}

		// write processed line to output
		output.WriteString(line + "\n")
	}

	return
}

func swap(oldfile, newfile string) (err error) {

	swapfile := oldfile + ".swap"

	err = os.Rename(newfile, swapfile)
	if err != nil {
		return
	}

	err = os.Rename(oldfile, newfile)
	if err != nil {
		return
	}

	err = os.Rename(swapfile, oldfile)
	if err != nil {
		return
	}

	return
}
