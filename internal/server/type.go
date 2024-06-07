package server

import "errors"

var (
	UnmarshallRequestError = errors.New("unmarshall request error")
	MarshallResponseError  = errors.New("marshall response error")
)
