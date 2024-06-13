package live_debugger

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

// LogCreatedDTO - server.InHandler outgoing dto
type LogCreatedDTO struct {
	Offset *uint64 `json:"offset"`
	Error  *string `json:"error"`
}

// LogDTO - server.InHandler incoming dto
type LogDTO struct {
	Level       Level       `json:"level"`
	Type        Type        `json:"type"`
	Message     Message     `json:"message"`
	Source      Source      `json:"source"`
	Fingerprint Fingerprint `json:"fingerprint"`
}
