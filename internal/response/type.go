package response

type LogResponse struct {
	Offset *uint64 `json:"offset"`
	Error  *string `json:"error"`
}
