package utils

import "strings"

// path: a/b/c/d
// output: a/b/c/
func GetParentPath(input string) string {
	lastSlashIndex := strings.LastIndex(input[:len(input)-1], "/")
	return input[:lastSlashIndex+1]
}

// path: a/b/c/d/
// output: d/
func GetChildPath(input string) string {
	lastSlashIndex := strings.LastIndex(input[:len(input)-1], "/")
	return input[lastSlashIndex+1:]
}

func RemoveFromSlice(slice []string, r string) []string {
	for i, v := range slice {
		if v == r {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
