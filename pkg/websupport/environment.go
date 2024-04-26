package websupport

import (
	"os"
	"strconv"
)

type intOrString interface {
	int | string
}

func EnvironmentVariable[T intOrString](variableName string, defaultValue T) T {
	value, found := os.LookupEnv(variableName)
	if !found {
		return defaultValue
	}

	var result T

	switch typedReference := any(&result).(type) {
	case *string:
		*typedReference = value
	case *int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		*typedReference = i
	}

	return result
}

func RequireEnvironmentVariable[T intOrString](variableName string) T {
	value, found := os.LookupEnv(variableName)
	if !found {
		panic("environment variable " + variableName + " is required")
	}

	var result T

	switch typedReference := any(&result).(type) {
	case *string:
		*typedReference = value
	case *int:
		i, err := strconv.Atoi(value)
		if err != nil {
			panic("environment variable " + variableName + " must be an integer")
		}
		*typedReference = i
	}

	return result
}
