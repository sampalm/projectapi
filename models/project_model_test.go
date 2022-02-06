package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/validator.v2"
)

func TestProjectModelFail(t *testing.T) {
	p := Project{}
	err := ""

	if validator.Validate(p) != nil {
		err = "invalid object"
	}

	assert.Equal(t, "invalid object", err)
}

func TestProjectModelValidationNameInvalidFail(t *testing.T) {
	p := Project{
		Name:        "invalid name $$",
		DisplayName: "Valid Display",
	}

	err := validator.Validate(p)

	assert.Equal(t, "Name: regular expression mismatch", err.Error())
}

func TestProjectModelValidationDisplayNameInvalidFail(t *testing.T) {
	p := Project{
		Name:        "validname",
		DisplayName: "va",
	}

	err := validator.Validate(p)

	assert.Equal(t, "DisplayName: less than min", err.Error())
}

func TestProjectModelValidationPass(t *testing.T) {
	p := Project{
		Name:        "validname",
		DisplayName: "Valid Display Name",
		Description: "valid description",
	}

	err := validator.Validate(p)

	assert.Equal(t, nil, err)
}
