package unique

// Ints returns a unique subset of the int slice provided.
func Ints(input []int) []int {
	u := make([]int, 0, len(input))
	m := map[int]bool{}
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

// Ints returns a unique subset of the string slice provided.
func Strings(input []string) []string {
	u := make([]string, 0, len(input))
	m := map[string]bool{}
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
