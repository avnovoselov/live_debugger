package response

import "encoding/json"

func EncodeJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}
