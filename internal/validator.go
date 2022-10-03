package internal

import "strconv"

type Validator interface {
	ValidInput(input string) (bool, int)
	TerminationInput(input string) bool
}

type numberValidator struct {
	terminationWord string
	requiredLen     int
}

func NewNumberValidator(terminationWord string, lenLimit int) Validator {
	return numberValidator{
		terminationWord: terminationWord,
		requiredLen:     lenLimit,
	}
}

func (nv numberValidator) ValidInput(input string) (bool, int) {
	if len(input) != nv.requiredLen {
		return false, 0
	}
	val, err := strconv.Atoi(input)
	if err != nil {
		return false, 0
	}
	return true, val
}

func (nv numberValidator) TerminationInput(input string) bool {
	return nv.terminationWord == input
}
