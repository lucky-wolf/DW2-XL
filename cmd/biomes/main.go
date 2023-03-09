package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

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

	// build the new map of race id -> factors
	header, races, err := load()
	if err != nil {
		return
	}

	// process our input + data -> output
	err = merge(header, races)
	if err != nil {
		return
	}

	// finally, swap the new races files so it's ready to go
	err = swap(oldfile, newfile)
	if err != nil {
		return
	}

	return
}

// we are meant to run from this folder
const datafile = "../../XL Biomes.csv"
const oldfile = "../../XL/Races.xml"
const newfile = "../../temp/Races.xml"
const swapfile = "../../XL/Races [new].xml"

// xml markers & templates
const start = `<ColonizationSuitabilityModifiers>`
const body = `			<OrbTypeFactor>
				<OrbTypeId>%s</OrbTypeId>
				<Factor>%s</Factor>
			</OrbTypeFactor>
`
const stop = `</ColonizationSuitabilityModifiers>`

var raceId = regexp.MustCompile(`\s+<RaceId>(\d+)</RaceId>`)

func load() (header []string, races map[string][]string, err error) {

	file, err := os.Open(datafile)
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
	header = strings.Split(scanner.Text(), ",")

	// maps racial id -> factors
	races = map[string][]string{}

	// process this line by line
	for line := 2; scanner.Scan(); line++ {

		// each line should contain the race name followed by the colonization factors
		values := strings.Split(scanner.Text(), ",")
		if len(values) == 1 {
			continue
		} else if len(values) != len(header) {
			err = fmt.Errorf("invalid line %d: not as many data elements as specified in header", line)
			return
		}

		races[values[0]] = values[1:]
	}
	return
}

func merge(header []string, races map[string][]string) (err error) {

	// open our races.xml as our input
	input, err := os.Open(oldfile)
	if err != nil {
		return
	}
	defer input.Close()

	output, err := os.Create(newfile)
	if err != nil {
		return
	}
	defer output.Close()

	// create a buffered writer
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// scan line by line
	scanner := bufio.NewScanner(input)

	currentRaceId := ""
	for lineIndex := 1; scanner.Scan(); lineIndex++ {

		// scan for our start token
		line := scanner.Text() + "\n"
		writer.WriteString(line)

		// skip everything until we get to the start of a new race block
		if strings.TrimSpace(line) != "<Race>" {
			continue
		}

		// skip everything until we find our raceId for this block as a whole
		for ; scanner.Scan(); lineIndex++ {
			line = scanner.Text() + "\n"
			writer.WriteString(line)
			if raceId.MatchString(line) {
				break
			}
		}

		// warn if our logic is bad
		if currentRaceId != "" {
			log.Printf("we hit a 2nd race id w/o finding our start token at line %d", lineIndex)
		}

		// record our active race id
		// https://pkg.go.dev/regexp#Regexp.FindStringSubmatch
		currentRaceId = raceId.FindStringSubmatch(line)[1]

		if currentRaceId == "" {
			log.Printf("Wtf?! no race id from %s", line)
		}

		// copy input until we find the start of the substitution block
	substitution:
		for ; scanner.Scan(); lineIndex++ {
			line = scanner.Text() + "\n"
			writer.WriteString(line)
			switch strings.TrimSpace(line) {
			case start:
				// write out our colonization factors for this race
				values := races[currentRaceId]
				for col := 0; col < len(values); col++ {
					_, err = writer.WriteString(fmt.Sprintf(body, header[col+1], values[col]))
					if err != nil {
						return
					}
				}

				// scan for our stop token
				for ; scanner.Scan(); lineIndex++ {
					line = scanner.Text()
					if strings.TrimSpace(line) == stop {
						break
					}
				}

				// go ahead and copy the stop line
				writer.WriteString(line + "\n")

				// we no longer have an active race id
				currentRaceId = ""
				break substitution

			case "</Race>":
				// no substitution block ever found for this race
				log.Printf("%d: no substitution block found for race id = %s", lineIndex, currentRaceId)
				// we no longer have an active race id
				currentRaceId = ""
				break substitution
			}
		}

	}

	return
}

func swap(oldfile, newfile string) (err error) {

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
