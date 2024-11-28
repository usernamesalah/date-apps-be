package util

func StringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}
