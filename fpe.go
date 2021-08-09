package tcg

// fpeEncode() represents the current FPE functionality.
func fpeEncode(pattern, input string) string {
	if pattern != `(\d\d)-(\d\d)` {
		panic("Unsupported pattern: " + pattern)
	}

	if input != "12-34" {
		panic("Unsupported value: " + input)
	}

	return "56-78"
}

func fpeDecode(pattern, encoded string) string {
	if pattern != `(\d\d)-(\d\d)` {
		panic("Unsupported pattern: " + pattern)
	}

	if encoded != "56-78" {
		panic("Unsupported value: " + encoded)
	}

	return "12-34"
}
