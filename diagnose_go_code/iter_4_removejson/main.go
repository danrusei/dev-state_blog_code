// Read the txt file
// Extract the ID (first column) and the JSON (the second column) that contain countries and continents.
// Display how many ID's are qualify as Luhn numbers
// Display the country with the highest number of hits
// Display the continent that has largest number of countries with hits over 1000

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/danrusei/dev-state_blog_code/tree/master/diagnose_go_code/iter_4_removejson/luhn"
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
	routines := runtime.NumCPU()
	results.getStatistics(file, routines)

	fmt.Printf("There are %d, out of %d, valid Luhn numbers. \n", results.validLuhn, results.nrLines)
	fmt.Printf("%s has the biggest # of visitors, with %d of hits. \n", results.country, results.hitsPerCountry)
	fmt.Printf("%s is the continent with most unique countries that accessed the site more than 1000 times. It has %d unique countries. \n", results.continent, results.countriesWithHits)

}

func (r *result) getStatistics(stream io.Reader, routines int) {

	countries := map[string]int{}
	continents := map[string]string{}

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	const CacheSize = 64 * 1024

	linesPool := sync.Pool{
		New: func() interface{} {
			return new([CacheSize]string)
		},
	}

	lines := make(chan *[CacheSize]string, routines)

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for cache := range lines {

				validLuhn := 0
				continent := ""
				country := ""

				for _, text := range cache {

					split := strings.Split(text, "#")
					if len(split) < 2 {
						continue
					}
					number := strings.TrimSpace(split[0])
					description := strings.TrimSpace(split[1])

					dep := strings.SplitN(description, "\"", 13)
					continent = dep[3]
					country = dep[11]

					if luhn.Valid(number) {
						validLuhn++
					}

					mutex.Lock()
					countries[country]++

					if _, ok := continents[country]; !ok {
						continents[country] = continent
					}
					mutex.Unlock()
				}

				mutex.Lock()
				r.validLuhn += validLuhn
				mutex.Unlock()

				linesPool.Put(cache)
			}
		}()
	}

	scanner := bufio.NewScanner(stream)
	iter := 0
	pool := linesPool.Get().(*[CacheSize]string)

	for {
		valid := scanner.Scan()
		if iter == CacheSize || !valid {
			lines <- pool
			iter = 0
			pool = linesPool.Get().(*[CacheSize]string)

		}
		if !valid {
			break
		}
		r.nrLines++
		pool[iter] = scanner.Text()
		iter++
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
