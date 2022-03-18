package cmd

const version = "v0.0.1"

func control(exp bool, truthValue, falsyValue string) string {
	if exp {
		return truthValue
	}

	return falsyValue
}
