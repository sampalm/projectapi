package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/validator.v2"
)

func TestApiModelValidationFail(t *testing.T) {
	a := API{}
	err := ""

	if validator.Validate(&a) != nil {
		err = "invalid object"
	}

	assert.Equal(t, "invalid object", err)

}

func TestApiModelValidationVersionFail(t *testing.T) {
	a := API{
		ProjectName: "thisprojectistest",
		Name:        "thisnameistest",
		Version:     "",
		OpenApiFile: "somebase64",
	}

	err := validator.Validate(&a)

	assert.Equal(t, "Version: zero value", err.Error())

}

func TestApiModelValidationNameMinLengthFail(t *testing.T) {
	a := API{
		ProjectName: "thisprojectistest",
		Name:        "12",
		Version:     "1.0",
		OpenApiFile: "somebase64",
	}

	err := validator.Validate(&a)

	assert.Equal(t, "Name: less than min", err.Error())

}

func TestApiModelValidationNameInvalidFail(t *testing.T) {
	a := API{
		ProjectName: "thisprojectistest",
		Name:        "this is not valid",
		Version:     "1.0",
		OpenApiFile: "somebase64",
	}

	err := validator.Validate(&a)

	assert.Equal(t, "Name: regular expression mismatch", err.Error())

}

func TestApiModelValidationProjectNameInvalidFail(t *testing.T) {
	a := API{
		ProjectName: "thisprojectis_not-valid#$%",
		Name:        "thisisvalid",
		Version:     "1.0",
		OpenApiFile: "somebase64",
	}

	err := validator.Validate(&a)

	assert.Equal(t, "ProjectName: regular expression mismatch", err.Error())

}

func TestApiModelValidationPass(t *testing.T) {
	a := API{
		ProjectName: "thisprojectisvalid",
		Name:        "thisisvalid",
		Version:     "1.0",
		OpenApiFile: "somebase64...",
	}

	err := validator.Validate(&a)

	assert.Equal(t, nil, err)

}
