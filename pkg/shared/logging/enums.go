package logging

import (
	"fmt"
	"log/slog"
	"strings"
)

type LogType int

const (
	LogTypeText LogType = iota
	LogTypeJSON
)

var logTypeFromString = map[string]LogType{
	"TEXT": LogTypeText,
	"JSON": LogTypeJSON,
}

func (logType LogType) String() string {
	return [...]string{"TEXT", "JSON"}[logType]
}

func ParseLogType(s string) (LogType, error) {
	normalized := strings.ToUpper(strings.TrimSpace(s))
	if val, ok := logTypeFromString[normalized]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("unknown log type: %q", s)
}

func ParseLogLevel(s string) (slog.Level, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown log level: %q", s)
	}
}
