package helper

import "strconv"

func GetParamAsInt32(key string) (int32, error) {
	intVal, err := strconv.Atoi(key)
	return int32(intVal), err
}
