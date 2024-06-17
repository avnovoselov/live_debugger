package configuration

import (
	"github.com/rs/zerolog"
)

// Common - root of configuration
type Common struct {
	Server Server `toml:"server"`
	Queue  Queue  `toml:"queue"`
	Logger Logger `toml:"logger"`
}

// Server - websocket server configuration
type Server struct {
	IP                               string `toml:"ip"`
	Port                             string `toml:"port"`
	InLocation                       string `toml:"in_location"`
	OutLocation                      string `toml:"out_location"`
	AmountOfSequentiallyErrorToBreak int    `toml:"amount_of_sequentially_error_to_break"`
	ThrottleDurationMs               int    `toml:"throttle_duration_ms"`
	SleepAfterErrorDurationMs        int    `toml:"sleep_after_error_duration_ms"`
}

// Queue - queue configuration
type Queue struct {
	Size uint64 `toml:"size"`
}

// Logger - live_debugger configuration
type Logger struct {
	Level zerolog.Level `toml:"level"`
}
