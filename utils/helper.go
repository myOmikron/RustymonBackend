package utils

func Contains[T comparable](i T, slice []T) bool {
	for _, elem := range slice {
		if elem == i {
			return true
		}
	}
	return false
}
