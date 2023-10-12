package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{map[string]string{}}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Check(cond bool, field, msg string) {
	if !cond {
		v.add(field, msg)
	}
}

func (v *Validator) add(field, msg string) {
	if _, ok := v.Errors[field]; !ok {
		v.Errors[field] = msg
	}
}
