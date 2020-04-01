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
	"runtime"
	"strings"
	"sync"

	"github.com/Danr17/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines/luhn"
)

//Region struct is used to Unmarshal the extracted JSON

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

	file, err := os.Open("../csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	routines := runtime.NumCPU() + 1

	validLuhn, nr, hitsCountry, hits, hitsContinent, hitsR := getStatistics(file, routines)

	fmt.Printf("There are %d, out of %d, valid Luhn numbers. \n", validLuhn, nr)
	fmt.Printf("%s has the biggest # of visitors, with %d of hits. \n", hitsCountry, hits)
	fmt.Printf("%s is the continent with most unique countries that accessed the site more than 1000 times. It has %d unique countries. \n", hitsContinent, hitsR)

}

func getStatistics(stream io.Reader, routines int) (int, int, string, int, string, int) {

	validLuhn := 0
	nr := 0
	countries := map[string]int{}
	continents := map[string]string{}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	lines := make(chan string, 2*routines)

	type region struct {
		Continent string `json:"continent"`
		Country   string `json:"country"`
	}

	for i := 0; i < 2*routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for text := range lines {

				split := strings.Split(text, "#")
				number := strings.TrimSpace(split[0])
				description := strings.TrimSpace(split[1])

				if luhn.Valid(number) {
					validLuhn++
				}

				var reg region

				err := json.Unmarshal([]byte(description), &reg)
				if err != nil {
					log.Println(err)
				}

				mutex.Lock()
				countries[reg.Country]++

				if _, ok := continents[reg.Country]; !ok {
					continents[reg.Country] = reg.Continent
				}
				mutex.Unlock()
			}
		}()
	}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		nr++
		lines <- scanner.Text()
	}

	close(lines)
	wg.Wait()

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
	return validLuhn, nr, hitsCountry, hits, hitsContinent, hitsR
}
