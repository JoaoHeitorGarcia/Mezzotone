package services

type SharedVariables map[string]any

var shared = NewSharedVariables()

func Shared() SharedVariables {
	return shared
}

func NewSharedVariables() SharedVariables {
	return make(SharedVariables)
}

func (sv SharedVariables) Get(key string) (any, bool) {
	v, ok := sv[key]
	return v, ok
}

func (sv SharedVariables) Set(key string, value any) {
	sv[key] = value
}

func (sv SharedVariables) Delete(key string) {
	delete(sv, key)
}
