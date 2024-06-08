package dto

import "encoding/json"

// EncodeJSON - marshal DTO to json bytes
func EncodeJSON[DTO any](v DTO) ([]byte, error) {
	return json.Marshal(v)
}

// ParseJSON - unmarshal json bytes to DTO
func ParseJSON[DTO any](b []byte) (dto DTO, err error) {
	err = json.Unmarshal(b, &dto)

	return
}
