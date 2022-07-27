package utils

type Errors struct {
	errors map[string][]any
}

func NewErrors() *Errors {
	err := &Errors{}
	err.errors = make(map[string][]any)

	return err
}

func (e *Errors) Get(key string) any {
	results := []any{}

	results, ok := e.errors[key]
	if ok {
		delete(e.errors, key)
	}

	return results
}

func (e *Errors) Set(key string, val any) {
	if dest, ok := e.errors[key]; ok {
		e.errors[key] = append(dest, val)
	} else {
		e.errors[key] = []any{val}
	}
}
