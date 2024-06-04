package request

type (
	Level       string
	Type        string
	Message     string
	Source      string
	Fingerprint string
)

const (
	Info    Level = "INFO"
	Warning Level = "WARNING"
	Error   Level = "ERROR"
)

type Request struct {
	Level       Level       `json:"level"`
	Type        Type        `json:"type"`
	Message     Message     `json:"message"`
	Source      Source      `json:"source"`
	Fingerprint Fingerprint `json:"fingerprint"`
}
