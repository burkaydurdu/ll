package cmd

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func control(exp bool, truthValue, falsyValue string) string {
	if exp {
		return truthValue
	}

	return falsyValue
}
