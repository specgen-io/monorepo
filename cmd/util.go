package cmd

func contains(array []string, that string) bool {
	for _, item := range array {
		if item == that {
			return true
		}
	}
	return false
}
