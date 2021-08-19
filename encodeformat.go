package tcg

import (
	"fmt"
	"regexp"
	"strings"
)

var debugging = false


func encodeFormatEncode(pattern, encode_format string, capture_groups []int, input string) string {
	re := regexp.MustCompile(pattern)

	allMatchedIndices := re.FindAllStringSubmatchIndex(input, -1)
	allMatches := re.FindAllStringSubmatch(input, -1)
	if len(allMatchedIndices) != 1 {
		panic(fmt.Sprintf("Expected a single match of pattern %s, but got %d for input %s", pattern, len(allMatchedIndices), input))
	}

	// [matchStart matchLen group1Start group1Len etc..]
	indices := allMatchedIndices[0]
	encodedMatches := fpeEncodeGroups(capture_groups, allMatches[0])

	debug("pattern: %s, encode_format: %s, capture_groups: %v, indices: %v, encodedMatches: %v\n",
		pattern, encode_format, capture_groups, indices, strings.Join(encodedMatches, ","))
	debug("Matched groups: %v\n", dump(re.FindAllStringSubmatch(input, -1)[0]))

	encodedInput := replaceGroups(input, encodedMatches, indices, capture_groups)
	debug("e input: %s\n-------- %s\n", input, encodedInput)

	result := []byte{}
	// result = re.Expand(result,[]byte(encode_format), []byte(encodedInput.String()), indices)
	result = re.Expand(result,[]byte(encode_format), []byte(encodedInput), indices)
	return string(result)
}

func fpeEncodeGroups(groupNumbers []int, groups []string) []string {
	result := make([]string, len(groups))
	copy(result, groups)

	for _, number := range groupNumbers {
		result[number] = strings.ToUpper(result[number])
	}

	return result
}

func dump(groups []string) string {
	var result strings.Builder
	for i, s := range groups {
		result.WriteString(fmt.Sprintf("%d:%s, ", i, s))
	}
	return result.String()
}

func encodeFormatDecode(pattern string, capture_groups []int, encodedInput string) string {
	debug("Decoding pattern: %s, capture_groups: %v, encodedInput: %s\n", pattern, capture_groups, encodedInput)
	re := regexp.MustCompile(pattern)

	allMatchedIndices := re.FindAllStringSubmatchIndex(encodedInput, -1)
	allMatches := re.FindAllStringSubmatch(encodedInput, -1)
	if len(allMatchedIndices) != 1 {
		panic(fmt.Sprintf("Expected a single match of pattern %s, but got %d for input %s", pattern, len(allMatchedIndices), encodedInput))
	}

	// [matchStart matchLen group1Start group1Len etc..]
	indices := allMatchedIndices[0]
	encodedMatches := fpeDecodeGroups(capture_groups, allMatches[0])

	debug("Decoding pattern: %s, capture_groups: %v, indices: %v, encodedMatches: %v\n",
		pattern, capture_groups, indices, dump(encodedMatches))

	decodedInput := replaceGroups(encodedInput, encodedMatches, indices, capture_groups)
	debug("d input: %s\n-------- %s\n", encodedInput, decodedInput)

	return decodedInput
}

func fpeDecodeGroups(groupNumbers []int, groups []string) []string {
	result := make([]string, len(groups))
	copy(result, groups)

	for _, number := range groupNumbers {
		result[number] = strings.ToLower(result[number])
	}

	return result
}

func debug(format string, args ...interface{}) {
	if debugging {
		fmt.Printf(format, args...)
	}
}

func replaceGroups(source string, groupMatches []string, indices []int, groups []int) string {
	result := source
	for _, group := range groups {
		groupIndex := 2 * group
		start, end := indices[groupIndex], indices[groupIndex+1]
		if start >= 0 {
			result = result[0:start] + groupMatches[group] + result[end:]
		}
	}

	return result
}
