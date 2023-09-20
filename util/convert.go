package util

func StringToInt(s string) int {
	var result int
	for _, v := range s {
		result = result*10 + int(v-'0')
	}
	return result
}
