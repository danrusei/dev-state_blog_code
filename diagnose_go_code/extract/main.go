//It's an ETL extracts from CSV only column 4 and 5, that is in our interest
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	csvfile, err := os.Open("../csv_files/test.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	txtfile, err := os.OpenFile("../csv_files/bench_test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer txtfile.Close()

	i, s, n := 0, 0, 0
	ID := ""

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		/*
		   //used to create bench file
		   	if n == 10000 {
		   		break
		   	}
		*/

		ID = record[4]

		switch {
		case i == 40:
			split := strings.SplitN(record[4], "", 4)
			ID = split[0] + " " + split[1] + " " + split[2] + " " + split[3]
			i = 0
		case s == 143:
			ID = "2334" + "$" + "1126677" + "notvalid"
			s = 0
		}

		txtfile.WriteString(ID + "#" + record[5] + "\n")
		i++
		s++
		n++
	}
}
