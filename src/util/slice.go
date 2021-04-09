package util

func SliceContainsString(s []string, e string) bool {
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
