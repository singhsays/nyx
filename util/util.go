package util

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	amountPattern = regexp.MustCompile("[()$,]")
)

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
