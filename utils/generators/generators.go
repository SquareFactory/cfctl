package generators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var matchCommaOutsideOfBrackets = regexp.MustCompile(`(?:\[+[^\[\]]*\]+|[^,])+`)

// SplitCommaOutsideOfBrackets generates strings
func SplitCommaOutsideOfBrackets(pattern string) []string {
	res := matchCommaOutsideOfBrackets.FindAllString(pattern, -1)
	if res == nil {
		return []string{""}
	}
	return res
}

// ExpandBrackets generates strings based on brackets ranges or digits.
//
// cn[1,2-4] generates cn1, cn2, cn4 and cn3.
func ExpandBrackets(pattern string) []string {
	var out []string
	if pattern == "" {
		return []string{}
	}

	beginIdx := -1
	for idx, rune := range pattern {
		// Search for '['
		if rune == '[' {
			beginIdx = idx
		}

		// Search for ']'
		if rune == ']' && beginIdx != -1 && beginIdx <= idx {
			digits := ParseRangeList(pattern[beginIdx+1 : idx])

			// Add the generated name
			for _, digit := range digits {
				out = append(out, fmt.Sprintf("%s%d%s", pattern[:beginIdx], digit, pattern[idx+1:]))
			}

			break
		}
	}

	// This means there is no brackets
	if beginIdx == -1 {
		return []string{pattern}
	}

	var merge []string
	for _, pattern := range out {
		names := ExpandBrackets(pattern)
		if len(names) == 0 {
			// Pattern is at the smallest factor
			merge = append(merge, pattern)
		} else {
			merge = append(merge, ExpandBrackets(pattern)...)
		}
	}

	return merge
}

// ParseRangeList converts a string containing comma separated digits and ranges into an array of digits.
//
// For example, "1,2-4" is [1,2,3,4].
func ParseRangeList(ranges string) []int {
	var digits []int
	digitsOrRanges := strings.Split(ranges, ",")

	for _, digitOrRange := range digitsOrRanges {
		// Check for digit
		digit, err := strconv.Atoi(digitOrRange)
		if err == nil {
			digits = append(digits, digit)
			continue
		}

		// Check for range
		r := strings.Split(digitOrRange, "-")
		if len(r) != 2 {
			// Is not a range
			continue
		}
		begin, errBegin := strconv.Atoi(r[0])
		end, errEnd := strconv.Atoi(r[1])
		if errBegin == nil && errEnd == nil {
			for digit := begin; digit <= end; digit++ {
				digits = append(digits, digit)
			}
			continue
		}
	}
	return digits
}
