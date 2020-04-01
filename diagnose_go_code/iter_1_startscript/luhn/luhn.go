package luhn

import (
	"log"
	"strconv"
	"unicode"
)

//Valid validates the string based on Luhn algorithm
func Valid(input string) bool {
	cleanInput := []int{}
	for _, s := range input {
		if unicode.IsDigit(s) {
			d, err := strconv.Atoi(string(s))
			if err != nil {
				log.Fatal(err)
			}
			cleanInput = append(cleanInput, d)
			continue
		}
		if s != ' ' {
			return false
		}
	}

	if len(cleanInput) < 2 {
		return false
	}

	sum := 0
	isEven := len(cleanInput)%2 == 0
	for _, val := range cleanInput {
		if isEven {
			val *= 2
			if val > 9 {
				val -= 9
			}
		}
		sum += val
		isEven = !isEven
	}

	return sum%10 == 0
}
