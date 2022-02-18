package flags

func Validate(flags ...Flag) error {
	for _, flag := range flags {
		if err := flag.Validate(); err != nil {
			return err
		}
	}
	return nil
}
