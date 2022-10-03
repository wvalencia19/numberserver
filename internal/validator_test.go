package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const terminationWord = "terminate"
const lenLimit = 9

func TestNumberValidator_ValidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid input",
			input: "123456789",
			want:  true,
		},
		{
			name:  "invalid length",
			input: "12345678",
			want:  false,
		},
		{
			name:  "not number with length",
			input: "1a3456789",
			want:  false,
		},
	}
	validator := NewNumberValidator(terminationWord, lenLimit)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := validator.ValidInput(tt.input)
			assert.Equal(t, valid, tt.want)
		})
	}
}

func TestNumberValidator_TerminationInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid termination word",
			input: terminationWord,
			want:  true,
		},
		{
			name:  "invalid termination word",
			input: "any",
			want:  false,
		},
	}
	validator := NewNumberValidator(terminationWord, lenLimit)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validator.TerminationInput(tt.input)
			assert.Equal(t, valid, tt.want)
		})
	}
}
