package tcg

import (
	"fmt"
	"regexp"
	"strings"
)

func captureGroupsEncode(pre_pattern string, capture_groups []int, pattern, input string) string {
	matchedPrePatternCaptureGroups := matchGroups(pre_pattern, capture_groups, input)
	fpeInput := computeFpeInput(matchedPrePatternCaptureGroups, pattern)

	return fpeEncode(pattern, fpeInput)

	//fmt.Printf("pre_pattern %s, capture_groups: %v, pattern: %s, input: %s, matched groups:%v, fpeInput: %s\n",
	//	pre_pattern, capture_groups, pattern, input, matchedPrePatternCaptureGroups, fpeInput )
}

// computeFpeInput constructs a string that when matched to 'pattern' the resulting matching groups
// are the same as 'groupValues'.
func computeFpeInput(groupValues []string, pattern string) string {
	patternGroupIndices := findGroupIndices(pattern)
	if len(groupValues) != len(patternGroupIndices) {
		panic("number of groups in 'pattern' must match that of 'capture_groups'")
	}

	var input strings.Builder
	patternIndex := 0
	for group, groupIndices := range patternGroupIndices {
		start, end := groupIndices[0], groupIndices[1]
		if patternIndex < start {
			// Copy everything before the opening paren
			input.WriteString(pattern[patternIndex:start])
		}
		// Copy the capture group value
		input.WriteString(groupValues[group])

		patternIndex = end;
	}
	if patternIndex < len(pattern) {
		input.WriteString(pattern[patternIndex:])
	}

	return input.String()
}

func findGroupIndices(pattern string) [][]int {
	// This is a toy implementation.
	// Is there a reliable way to implement this function?

	re := regexp.MustCompile(`(\([^\(]*\))`) // `(\(xxxx\))`) // where xxxx is [^\(]
	return re.FindAllStringIndex(pattern, -1)
}

// matchGroups returns a map where the key is the group number, and the value is the matched substring
// For example: matchGroups(`(\d\d)-(\d\d)-(\d\d)`, []int{1,3}, "12-34-56")
//     1: "12"
//     3: "56"
func matchGroups(pre_pattern string, capture_groups []int, input string) []string {
	preRegexp := regexp.MustCompile(pre_pattern)
	submatches := preRegexp.FindAllStringSubmatchIndex(input, -1)
	if len(submatches) != 1 {
		panic("Expected a single set of submatches")
	}

	var result []string
	for _, cg := range capture_groups {
		start := submatches[0][2 + 2*(cg-1)]
		end := submatches[0][3 + 2*(cg-1)]

		if start == -1 || end == -1 {
			panic(fmt.Sprintf("Required group not found: %d", cg))
		}
		result = append(result, string(input[start:end]))
	}

	return result
}
