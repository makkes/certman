package flags

import (
	"fmt"
	"strings"
)

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

type Validatable interface {
	Validate() error
}

type DisallowEmpty struct{}

func (v DisallowEmpty) Validate(val string) error {
	if val == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

type exclusive struct {
	flags []Flag
}

func Exclusive(flags ...Flag) exclusive {
	return exclusive{flags: flags}
}

func (e exclusive) Validate() error {
	hasNonEmpty := false
	for _, f := range e.flags {
		if f.Val != "" {
			if hasNonEmpty {
				return e.error()
			}
			hasNonEmpty = true
		}
	}
	return nil
}

func (e exclusive) error() error {
	names := make([]string, len(e.flags))
	for idx, f := range e.flags {
		names[idx] = "--" + f.Name
	}
	return fmt.Errorf("only one of %s must be set", strings.Join(names, ", "))
}

type and struct {
	flags []Flag
}

func And(flags ...Flag) and {
	return and{flags: flags}
}

func (a and) Validate() error {
	for _, f := range a.flags {
		if f.Val == "" {
			return a.error()
		}
	}
	return nil
}

func (a and) error() error {
	names := make([]string, len(a.flags))
	for idx, f := range a.flags {
		names[idx] = "--" + f.Name
	}
	return fmt.Errorf("all of %s must be set", strings.Join(names, ", "))
}
