package tid

func first[S ~[]E, E any](s S) (v E) {
	if len(s) > 0 {
		return s[0]
	}

	return
}
