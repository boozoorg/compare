package compare

// compare to comparable type
func IsEqual[K comparable](first, second K) bool {
	return first == second
}

// Todo:add compare of to not comparable type
