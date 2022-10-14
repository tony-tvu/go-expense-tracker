package util

func ContainsInt(arr *[]int, find int) bool {
	for _, el := range *arr {
		if find == el {
			return true
		}
	}
	return false
}
