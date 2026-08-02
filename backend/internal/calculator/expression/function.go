package expression

// IsSupportedFunction reports whether name is a known function.
// Names are expected to already be lowercased by the tokenizer.
func IsSupportedFunction(name string) bool {
	return name == "sqrt"
}
