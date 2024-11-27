package models

import (
	"time"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	LogToConsole bool   `yaml:"logToConsole"`
	LogToFile    bool   `yaml:"logToFile"`
	LogFilePath  string `yaml:"logFilePath"`
	MaxFileSize  int    `yaml:"maxFileSize"`
	MaxBackups   int    `yaml:"maxBackups"`
	MaxAge       int    `yaml:"maxAge"`
}

// ConsoleFormat 定义控制台输出格式的配置
type ConsoleFormat struct {
	TimeFormat   string
	NoColor      bool
	PartsOrder   []string
	PartsExclude []string
}

// DefaultConsoleFormat 返回默认的控制台格式配置
func DefaultConsoleFormat() ConsoleFormat {
	return ConsoleFormat{
		TimeFormat: time.DateTime,
		NoColor:    false,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
		},
	}
}
