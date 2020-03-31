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

	txtfile, err := os.OpenFile("../csv_files/test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer txtfile.Close()

	i, s := 0, 0
	ID := ""

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		ID = record[3]

		switch {
		case i == 40:
			split := strings.SplitN(record[3], "", 4)
			ID = split[0] + " " + split[1] + " " + split[2] + " " + split[3]
			i = 0
		case s == 143:
			ID = "2334" + "$" + "1126677" + "notvalid"
			s = 0
		}

		txtfile.WriteString(ID + "#" + record[4] + "\n")
		i++
		s++
	}
}
