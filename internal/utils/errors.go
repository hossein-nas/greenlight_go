package utils

func CheckError(val interface{}) bool {
	if _, ok := val.(error); ok {
		return true
	} else {
		return false
	}
}
