package util

// UniqInts returns a unique subset of the int slice provided.
func UniqInts(input []int) []int {
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

// UniqString returns a unique subset of the string slice provided.
func UniqStrings(input []string) []string {
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
