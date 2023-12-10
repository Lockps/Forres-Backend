package function

import "strings"

func StrToByteSlice(data string) []byte {
	return []byte(data)
}

func StrSliceToByteSlice(strSlice []string) []byte {
	joinedString := strings.Join(strSlice, ",")

	byteSlice := []byte(joinedString)

	return byteSlice
}

func ChecktheError(x error) error {
	if x != nil {
		return x
	}
	return nil
}
