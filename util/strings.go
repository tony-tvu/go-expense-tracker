package util

func ContainsEmpty(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return true
		}
	}
	return false
}

func Contains(arr *[]string, find string) bool {
	for _, word := range *arr {
		if find == word {
			return true
		}
	}
	return false
}