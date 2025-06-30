package util

import "slices"

func RemoveFirstString(slice []string, target string) []string {
	copied := append([]string(nil), slice...)
	for i, s := range slice {
		if s == target {
			return slices.Delete(copied, i, i+1)
		}
	}
	return slice
}
