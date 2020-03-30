// Read the txt file
// Apply the Luhn algorithm over the first column
// Extract the Json from the second column and populate the result into map[string]map[string]string :: Continent(Country(#))

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/pkg/luhn"
)

var validLuhn int64
var nr int64
var romania int64

type Region struct {
	Continent string `json:"continent"`
	Country   string `json:"country"`
}

func main() {

	file, err := os.Open("../../csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

		if reg.Country == "Romania" {
			romania++
		}

	}
	fmt.Printf("%d are valid Luhn numbers, out of %d. Romania was mentioned %d times", validLuhn, nr, romania)
}
