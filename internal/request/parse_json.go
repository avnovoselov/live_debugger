package request

import (
	"encoding/json"
)

func ParseJSON(b []byte) (req LogRequest, err error) {
	err = json.Unmarshal(b, &req)

	return
}
