package util

func RmStrFromAry(slice []string, target string) []string {
	for s := 0; s < len(slice); s++ {
		if slice[s] == target {
			return append(slice[:s], slice[s+1:]...)
		}
	}
	return slice
}

func RmDuplicatesStr(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if _, ok := encountered[elements[v]]; !ok {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
