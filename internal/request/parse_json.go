package request

import (
	"encoding/json"
)

func ParseJSON(b []byte) (LogRequest, error) {
	request := LogRequest{}
	err := json.Unmarshal(b, &request)

	return request, err
}
