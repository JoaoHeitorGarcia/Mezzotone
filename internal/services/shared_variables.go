package services

var sharedVariables = make(map[string]any)

func getSharedVariables() map[string]any {
	return sharedVariables
}

func getSharedVariable(key string) any {
	return sharedVariables[key]
}

func SetSharedVariable(key string, value string) {
	sharedVariables[key] = value
}

func DeleteSharedVariable(key string) {
	delete(sharedVariables, key)
}
