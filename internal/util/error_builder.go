package util

import "errors"

// ErrorBuilder - returns a function with the given preset base error.
//
// Example:
//
//	eb := ErrorBuilder(errors.New("foo"))
//	err := eb(errors.New("bar"))
//	fmt.Println(err)
//	> "foo\nbar"
func ErrorBuilder(baseError error) func(extendedError error) error {
	return func(extendedError error) error {
		return errors.Join(baseError, extendedError)
	}
}
