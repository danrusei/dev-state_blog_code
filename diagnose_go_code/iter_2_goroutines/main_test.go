package main

import (
	"log"
	"os"
	"runtime"
	"testing"
)

func BenchmarkGetStatistics(b *testing.B) {

	file, err := os.Open("../csv_files/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b.ResetTimer()

	results := result{}
	routines := runtime.NumCPU() * 2

	for i := 0; i < b.N; i++ {
		results.getStatistics(file, routines)
	}
}
