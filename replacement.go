package tcg

import "regexp"

func replacementEncode(pre_pattern, pre_replacement, pattern, input string) string {
	processedInput := preProcessInput(pre_pattern, pre_replacement, input)

	return fpeEncode(pattern, processedInput)
}

func preProcessInput(pre_pattern string, pre_replacement string, input string) string {
	preRegexp := regexp.MustCompile(pre_pattern)

	result := []byte{}
	submatches := preRegexp.FindAllStringSubmatchIndex(input, -1)
	if len(submatches) != 1 {
		panic("Expected a single set of submatches")
	}

	return string(preRegexp.ExpandString(result, pre_replacement, input, submatches[0]))
}
