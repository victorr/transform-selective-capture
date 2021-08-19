package tcg

import (
	"fmt"
	"testing"
)

func TestEncodeFormatEncode(t *testing.T) {
	testCases := []struct {
		pattern        string
		encodeFormat   string
		capture_groups []int
		input          string
		output         string
		decoded        string
	}{  // pattern                                    encodeFormat,      captureGroups, input,     output, decoded
		{`([[:alpha:]]{2})-([[:alpha:]]{2})`,              "$1-$2",      []int{1,2}, "ab-cd",     "AB-CD", "ab-cd"},  // reality check: no modification
		{`(ccn:|card:)?([[:alpha:]]{2})-([[:alpha:]]{2})`, "$2-$3",      []int{2,3}, "ccn:ab-cd", "AB-CD", "ab-cd"},  // Prefix
		{`(ccn:|card:)?([[:alpha:]]{2})-([[:alpha:]]{2})`, "$2-$3",      []int{2,3}, "card:ab-cd","AB-CD", "ab-cd"},
		{`(ccn:|card:)?([[:alpha:]]{2})-([[:alpha:]]{2})`, "card:$2-$3", []int{2,3}, "cnn:ab-cd", "card:AB-CD", "card:ab-cd"}, // reformat - force prefix to be card:
		{`(ccn:|card:)?([[:alpha:]]{2})-([[:alpha:]]{2})`, "$1$2-$3",    []int{2,3}, "card:ab-cd","card:AB-CD", "card:ab-cd"}, // pass through
		{`(ccn:|card:)?([[:alpha:]]{2})-([[:alpha:]]{2})`, "$2-$3",      []int{2,3}, "ab-cd",     "AB-CD", "ab-cd"},
		{`([[:alpha:]]{2})[\.\-/ ]([[:alpha:]]{2})`,       "$1-$2",      []int{1,2}, "ab.cd",     "AB-CD", "ab-cd"},   // Separator is a dot, dash, slash or space
		{`([[:alpha:]]{2})[\.\-/ ]([[:alpha:]]{2})`,       "$1-$2",      []int{1,2}, "ab/cd",     "AB-CD", "ab-cd"},   // Separator is a dot, dash, slash or space
		{`([[:alpha:]]{2})[\.\-/ ]([[:alpha:]]{2})`,       "$1-$2",      []int{1,2}, "ab-cd",     "AB-CD", "ab-cd"},   // Separator is a dot, dash, slash or space
		{`([[:alpha:]]{2})[\.\-/ ]([[:alpha:]]{2})`,       "$1-$2",      []int{1,2}, "ab cd",     "AB-CD", "ab-cd"},   // Separator is a dot, dash, slash or space
		{`([[:alpha:]]{2})<->([[:alpha:]]{2})|([[:alpha:]]{2})-([[:alpha:]]{2})`, "$2-$1", []int{1,2,3,4}, "cd<->ab",   "AB-CD", "ab-cd"}, // Reordering of the input    <<<<<<<<

		// phone number spaces to dashes 123 456 7890, leave the last four numbers in the clear
		// {`(...) (...) (....)`, "$1-$2-$3", []int{1,2}, "abc-def-ghij", "ABC-DEF-ghij", "abc-def-ghij", }, // not possible to decode it
		{`(...)[ -](...)[ -](....)`, "$1-$2-$3", []int{1,2}, "abc-def-ghij", "ABC-DEF-ghij", "abc-def-ghij", },
		{`(...)[ -](...)[ -](....)`, "$1-$2-$3", []int{1,2}, "abc def ghij", "ABC-DEF-ghij", "abc-def-ghij", },
		{`(...)[ -](...)[ -](....)`, "$1-$2-$3", []int{1,2}, "abc def ghij", "ABC-DEF-ghij", "abc-def-ghij", },

		// userid / email
		// {`((user_id:(.+))|(email:(.+@.+)))`, "$3$5", []int{3,5}, "user_id:the-userid", "THE-USERID", "the-userid",}, // FAILS, since THE-USER-ID has no prefix
		{`((user_id:)(.+)|(email:)(.+@.+)|(.+))`, "$3$5", []int{3,5,6}, "user_id:the-userid", "THE-USERID", "the-userid",}, //
		{`((user_id:)(.+)|(email:)(.+@.+)|(.+))`, "$3$5", []int{3,5,6}, "email:the@user.id", "THE@USER.ID", "the@user.id",}, //

		// SSN
		{`(SSN:)?(...)[\.-](..)[\.-](....)`, "$2-$3-$4", []int{2,3,4}, "aaa.bb.cccc", "AAA-BB-CCCC", "aaa-bb-cccc", },
		{`(SSN:)?(...)[\.-](..)[\.-](....)`, "$2-$3-$4", []int{2,3,4}, "aaa-bb-cccc", "AAA-BB-CCCC", "aaa-bb-cccc", },
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%s", i, tc.input), func(t *testing.T) {
			encoded := encodeFormatEncode(tc.pattern, tc.encodeFormat, tc.capture_groups, tc.input)

			if encoded != tc.output {
				t.Errorf("Encoding failed, expected %s, but got %s", tc.output, encoded)
			}

			decoded := encodeFormatDecode(tc.pattern, tc.capture_groups, encoded)
			if decoded != tc.decoded {
				t.Errorf("Decoding failed, expected %s, but got %s", tc.decoded, decoded)
			}
		} )
	}
}
