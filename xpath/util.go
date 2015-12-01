package xpath

// String returns the String component of the result, and as a side effect
// releases the Result by calling Free() on it
func String(r Result) string {
	defer r.Free()
	return r.String()
}

func Bool(r Result) bool {
	defer r.Free()
	return r.Bool()
}

func Number(r Result) float64 {
	defer r.Free()
	return r.Number()
}

