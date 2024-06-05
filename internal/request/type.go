package request

type (
	Level       int
	Type        string
	Message     string
	Source      string
	Fingerprint string
)

const (
	Debug   Level = iota
	Info    Level = iota
	Warning Level = iota
	Error   Level = iota
)

type LogRequest struct {
	Level       Level       `json:"level"`
	Type        Type        `json:"type"`
	Message     Message     `json:"message"`
	Source      Source      `json:"source"`
	Fingerprint Fingerprint `json:"fingerprint"`
}

type LogResponse struct {
	Offset *uint64 `json:"offset"`
	Error  *string `json:"error"`
}
