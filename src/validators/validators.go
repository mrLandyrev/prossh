package validators

type IValidator interface {
	Validate(input string) (err string)
}

type MinLengthValidator struct {
	Value int
}

func (validator MinLengthValidator) Validate(input string) (err string) {
	if len(input) < validator.Value {
		err = "Lenght is small"
	}
	return
}

type CannotExistsValidator struct {
	Items []string
}

func (validator CannotExistsValidator) Validate(input string) (err string) {
	for _, item := range validator.Items {
		if item == input {
			err = "Item exsists in array"
			return
		}
	}
	return
}

type ExistsValidator struct {
	Items []string
}

func (validator ExistsValidator) Validate(input string) (err string) {
	for _, item := range validator.Items {
		if item == input {
			return
		}
	}
	err = "Item not exsists in array"
	return
}
