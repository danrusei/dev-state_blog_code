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

	"github.com/danrusei/dev-state_blog_code/tree/master/diagnose_go_code/iter_2_goroutines/luhn"
)

type result struct {
	validLuhn         int
	nrLines           int
	country           string
	hitsPerCountry    int
	continent         string
	countriesWithHits int
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

	file, err := os.Open("../csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	results := result{}
	routines := runtime.NumCPU() * 2
	results.getStatistics(file, routines)

	fmt.Printf("There are %d, out of %d, valid Luhn numbers. \n", results.validLuhn, results.nrLines)
	fmt.Printf("%s has the biggest # of visitors, with %d of hits. \n", results.country, results.hitsPerCountry)
	fmt.Printf("%s is the continent with most unique countries that accessed the site more than 1000 times. It has %d unique countries. \n", results.continent, results.countriesWithHits)

}

func (r *result) getStatistics(stream io.Reader, routines int) {

	//region struct is used to Unmarshal the JSON
	type region struct {
		Continent string `json:"continent"`
		Country   string `json:"country"`
	}

	countries := map[string]int{}
	continents := map[string]string{}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	lines := make(chan string, routines)

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for text := range lines {

				split := strings.Split(text, "#")
				if len(split) < 2 {
					continue
				}
				number := strings.TrimSpace(split[0])
				description := strings.TrimSpace(split[1])

				if luhn.Valid(number) {
					mutex.Lock()
					r.validLuhn++
					mutex.Unlock()
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
		r.nrLines++
		lines <- scanner.Text()
	}

	close(lines)
	wg.Wait()

	for k, v := range countries {
		if v > r.hitsPerCountry {
			r.hitsPerCountry = v
			r.country = k
		}
	}

	regions := map[string]int{}
	for k, v := range continents {
		if countries[k] > 1000 {
			regions[v]++
		}
	}

	r.countriesWithHits = 1
	for k, v := range regions {
		if v > r.countriesWithHits {
			r.continent = k
			r.countriesWithHits = v
		}
	}
}
