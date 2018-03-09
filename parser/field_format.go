package parser

import (
	"fmt"
	"strings"

	"github.com/fatih/camelcase"
)

var defaultFieldFormat = map[string]FieldFormat{
	"json": CamelCase,
	"db":   SnakeCase,
}

// FieldFormat field name format
type FieldFormat string

// Available field format
const (
	Default   FieldFormat = ""
	CamelCase             = "camel"
	SnakeCase             = "snake"
)

func (f FieldFormat) validate() error {
	switch f {
	case Default, CamelCase, SnakeCase:
		return nil
	}

	return fmt.Errorf("invalid field format %q", f)
}

func (f FieldFormat) transform(key, value string) string {
	if f == Default {
		dfmt := defaultFieldFormat[key]
		if dfmt == Default {
			dfmt = CamelCase
		}

		return dfmt.transform(key, value)
	}

	splitted := camelcase.Split(value)
	parts := make([]string, len(splitted))

	var sep string
	var transFn func(string) string

	switch f {
	case CamelCase:
		transFn = strings.Title

	case SnakeCase:
		transFn = strings.ToLower
		sep = "_"
	}

	parts[0] = strings.ToLower(splitted[0])
	for i := 1; i < len(splitted); i++ {
		parts[i] = transFn(splitted[i])
	}

	return strings.Join(parts, sep)
}
