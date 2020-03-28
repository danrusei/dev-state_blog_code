package luhn

// Valid validates the giving string to be luhn compliant.
func Valid(s string) bool {
	var idx, sum int

	// need to go backwards to be able to work with 2nd digits from the right
	// as spaces could be involved along the way
	for i := len(s) - 1; i >= 0; i-- {
		j := rune(s[i])
		// handling spaces inside the loop, instead of strings.Replace is a huge difference
		if j == ' ' {
			continue
		}

		// bit faster as unicode.IsDigit()
		if j < '0' || j > '9' {
			return false
		}
		// keep track of actual index
		idx++

		// get the integer value
		n := int(j - '0')

		if idx%2 == 0 {
			// nerdy n *= 2
			n <<= 1
			if n > 9 {
				n -= 9
			}
		}

		sum += n
	}

	// due to spaces handled in the loop, now we can see
	// if we have enough digits
	if idx <= 1 {
		return false
	}

	return sum%10 == 0
}
