package main

import (
	"log"
	"os"
	"testing"
)

func BenchmarkGetStatistics(b *testing.B) {

	file, err := os.Open("csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getStatistics(file)
	}
}
