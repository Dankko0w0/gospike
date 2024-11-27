package logger

import (
	"io"
	"os"
	"sync"

	"github.com/Dankko0w0/gospike/models"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger zerolog.Logger
	once   sync.Once
)

// InitializeLogger initializes the logger with specified settings
func InitializeLogger(logToConsole bool, logToFile bool, logFilePath string, maxFileSize int, maxBackups int, maxAge int, consoleFormat *models.ConsoleFormat) {
	once.Do(func() {
		var writers []io.Writer

		if logToConsole {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

			// 如果提供了自定义格式，则应用它
			if consoleFormat != nil {
				consoleWriter.TimeFormat = consoleFormat.TimeFormat
				consoleWriter.NoColor = consoleFormat.NoColor
				if len(consoleFormat.PartsOrder) > 0 {
					consoleWriter.PartsOrder = consoleFormat.PartsOrder
				}
				if len(consoleFormat.PartsExclude) > 0 {
					consoleWriter.PartsExclude = consoleFormat.PartsExclude
				}
			}

			writers = append(writers, consoleWriter)
		}

		if logToFile {
			writers = append(writers, &lumberjack.Logger{
				Filename:   logFilePath,
				MaxSize:    maxFileSize, // megabytes
				MaxBackups: maxBackups,
				MaxAge:     maxAge, // days
			})
		}

		multi := io.MultiWriter(writers...)
		logger = zerolog.New(multi).With().Timestamp().Logger()
	})
}

// Info logs an info message
func Info(msg string) {
	logger.Info().Msg(msg)
}

// Error logs an error message
func Error(msg string, err error) {
	logger.Error().Err(err).Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	logger.Debug().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	logger.Warn().Msg(msg)
}
