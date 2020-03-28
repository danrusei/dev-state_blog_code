package luhn

import (
	"strings"
	"unicode"
)

//Valid checks if the number is a Luhn numbers
func Valid(input string) bool {

	input = strings.ReplaceAll(input, " ", "")

	if len(input) < 2 {
		return false
	}

	sum := 0
	isEven := len(input)%2 == 0
	for _, r := range input {

		if !unicode.IsDigit(r) {
			return false
		}
		val := int(r - '0')
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
