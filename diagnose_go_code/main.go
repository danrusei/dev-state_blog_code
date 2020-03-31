// Read the txt file
// Extract the ID (first column) and the JSON (the second column) that contain countries and continents.
// Display how many ID's are qualify as Luhn numbers
// Display the country with the highest number of hits
// Display the continent that has largest number of countries with hits over 1000

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/trace"
	"strings"

	"github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/luhn"
)

var validLuhn int64
var nr int64

//Region struct is used to Unmarshal the extracted JSON
type Region struct {
	Continent string `json:"continent"`
	Country   string `json:"country"`
}

func main() {

	/*
	tracefile, err := os.OpenFile("m.trace", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer tracefile.Close()

	trace.Start(tracefile)
	defer trace.Stop()
	*/

	file, err := os.Open("csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	countries, continents := getStatistics(file)

	hits := 0
	hitsCountry := ""
	for k, v := range countries {
		if v > hits {
			hits = v
			hitsCountry = k
		}
	}

	regions := map[string]int{}
	for k, v := range continents {
		if countries[k] > 1000 {
			regions[v]++
		}
	}

	hitsR := 1
	hitsContinent := ""
	for k, v := range regions {
		if v > hitsR {
			hitsContinent = k
			hitsR = v
		}
	}

	fmt.Printf("There are %d, out of %d, valid Luhn numbers. \n", validLuhn, nr)
	fmt.Printf("%s has the biggest # of visitors, with %d of hits. \n", hitsCountry, hits)
	fmt.Printf("%s is the continent with most unique countries that accessed the site more than 1000 times. It has %d unique countries. \n", hitsContinent, hitsR)

}

func getStatistics(stream io.Reader) (map[string]int, map[string]string) {
	countries := map[string]int{}
	continents := map[string]string{}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		nr++

		text := scanner.Text()

		split := strings.Split(text, "#")
		number := strings.TrimSpace(split[0])
		description := strings.TrimSpace(split[1])

		if luhn.Valid(number) {
			validLuhn++
		}

		var reg Region

		err := json.Unmarshal([]byte(description), &reg)
		if err != nil {
			log.Println(err)
		}

		countries[reg.Country]++

		if _, ok := continents[reg.Country]; !ok {
			continents[reg.Country] = reg.Continent
		}

	}
	return countries, continents
}
