package alfred

import "time"

func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

func FilterSlice[T any](data []T, fn func(v T) bool) []T {
	var result []T
	for _, v := range data {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func FilterMap[T any](data map[string]T, fn func(k string, v T) bool) map[string]T {
	result := make(map[string]T)
	for k, v := range data {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}
