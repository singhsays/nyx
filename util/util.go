package util

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	amountPattern = regexp.MustCompile("[()$,]")
)

// Threshold for floating point comparisons.
const THRESHOLD = 0.00000001

// ToAmount parses a currency string into a float.
// It
// 	- infers a negative amount from enclosing parentheses or a - prefix
//  - it removes commas and $ signs.
//  - parses the string value into a float64.
func ToAmount(value string) (float64, error) {
	multiplier := 1.0
	if strings.HasPrefix(value, "(") {
		multiplier = -1.0
	}
	amount, err := strconv.ParseFloat(amountPattern.ReplaceAllString(value, ""), 64)
	if err != nil {
		return 0, err
	}
	return multiplier * amount, nil
}

// FloatsEqual checks if the given floats are (nearly) equal.
func FloatsEqual(first, second float64) bool {
	return (math.Abs(second-first) < THRESHOLD)
}
