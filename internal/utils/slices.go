package utils

func StringSlicePairs(slice []string) map[string]string {
	if len(slice) == 0 {
		return map[string]string{}
	}

	result := map[string]string{}
	for i := 0; i < len(slice); i += 2 {
		if i+1 == len(slice) {
			result[slice[i]] = ""
		} else {
			result[slice[i]] = slice[i+1]
		}
	}

	return result
}
