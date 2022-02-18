package flags

import "fmt"

type Flag struct {
	Name      string
	Val       string
	Validator Validator
}

func (f Flag) Validate() error {
	if f.Validator == nil {
		return nil
	}
	if err := f.Validator.Validate(f.Val); err != nil {
		return fmt.Errorf("flag --%s is not valid: %w", f.Name, err)
	}
	return nil
}

type Validator interface {
	Validate(val string) error
}

type DisallowEmpty struct{}

func (v DisallowEmpty) Validate(val string) error {
	if val == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}
