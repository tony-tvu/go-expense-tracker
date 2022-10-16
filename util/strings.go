package util

import "regexp"

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

func RemoveDuplicateWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}