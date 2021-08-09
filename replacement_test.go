package tcg

import "testing"

func TestReplacementEncode(t *testing.T) {
	testCases := []struct {
		pre_pattern   string
		pre_placement string
		pattern       string
		input         string
	}{
		{`(\d\d)-(\d\d)`,               "$1-$2", `(\d\d)-(\d\d)`, "12-34"},   // reality check: no modification
		{ `(ccn:|card:)?(\d\d)-(\d\d)`, "$2-$3", `(\d\d)-(\d\d)`, "ccn:12-34" }, // Prefix
		{ `(ccn:|card:)?(\d\d)-(\d\d)`, "$2-$3", `(\d\d)-(\d\d)`, "card:12-34" },
		{ `(ccn:|card:)?(\d\d)-(\d\d)`, "$2-$3", `(\d\d)-(\d\d)`, "12-34" },
		{`(\d\d)[\.\-/ ](\d\d)`,        "$1-$2", `(\d\d)-(\d\d)`, "12.34"},   // Separator is a dot, dash, slash or space
		{`(\d\d)[\.\-/ ](\d\d)`,        "$1-$2", `(\d\d)-(\d\d)`, "12/34"},   // Separator is a dot, dash, slash or space
		{`(\d\d)[\.\-/ ](\d\d)`,        "$1-$2", `(\d\d)-(\d\d)`, "12-34"},   // Separator is a dot, dash, slash or space
		{`(\d\d)[\.\-/ ](\d\d)`,        "$1-$2", `(\d\d)-(\d\d)`, "12 34"},   // Separator is a dot, dash, slash or space
		{`(\d\d)<->(\d\d)`,             "$2-$1", `(\d\d)-(\d\d)`, "34<->12"}, // Reordering of the input
	}

	for _, tc := range testCases {
		encoded := replacementEncode(tc.pre_pattern, tc.pre_placement, tc.pattern, tc.input)

		if encoded != "56-78" {
			t.Errorf("Expected 56-78, but got %s", encoded)
		}
	}
}
