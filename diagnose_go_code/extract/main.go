package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func main() {

	csvfile, err := os.Open("../csv_files/test.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)

	txtfile, err := os.OpenFile("../csv_files/test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer txtfile.Close()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		txtfile.WriteString(record[3] + "#" + record[4] + "\n")
	}
}
