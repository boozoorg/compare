package compare

func IsEqual[K comparable](first, second K) bool {
	return first == second
}
