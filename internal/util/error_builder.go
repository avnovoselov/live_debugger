package util

import "errors"

func ErrorBuilder(baseError error) func(extendedError error) error {
	return func(extendedError error) error {
		return errors.Join(baseError, extendedError)
	}
}
