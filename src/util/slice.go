package util

type Comparable interface {
	string | rune
}

func SliceContains[T Comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicates(arr []string) []string {
	set := map[string]struct{}{}

	for _, s := range arr {
		set[s] = struct{}{}
	}
	res := []string{}
	for lang := range set {
		res = append(res, lang)
	}
	return res
}
