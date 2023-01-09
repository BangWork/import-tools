package utils

// StringArrayDifference additions = new - old, deletions = old - new
func StringArrayDifference(old []string, new []string) (additions []string, deletions []string) {
	additionsMap := make(map[string]struct{})
	deletionsMap := make(map[string]struct{})
	for _, s := range old {
		if len(s) > 0 {
			deletionsMap[s] = struct{}{}
		}
	}
	for _, s := range new {
		if len(s) > 0 {
			additionsMap[s] = struct{}{}
		}
	}
	for s, _ := range additionsMap {
		if _, ok := deletionsMap[s]; !ok {
			additions = append(additions, s)
		}
	}
	for s, _ := range deletionsMap {
		if _, ok := additionsMap[s]; !ok {
			deletions = append(deletions, s)
		}
	}
	return
}

func UniqueNoNullSlice(slice ...string) (newSlice []string) {
	found := make(map[string]bool)
	for _, val := range slice {
		if val == "" {
			continue
		}
		if _, ok := found[val]; !ok {
			found[val] = true
			newSlice = append(newSlice, val)
		}
	}
	return
}

func StringArrayIntersection(a []string, b []string) []string {
	aset := make(map[string]struct{})
	for _, s := range a {
		aset[s] = struct{}{}
	}
	rset := make(map[string]struct{})
	for _, s := range b {
		if _, ok := aset[s]; ok {
			rset[s] = struct{}{}
		}
	}
	result := make([]string, len(rset))
	i := 0
	for s, _ := range rset {
		result[i] = s
		i++
	}
	return result
}

func StringArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
