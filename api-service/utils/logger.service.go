package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath" // For getting base filename
	"runtime"
	"strconv"
	"strings" // For cleaning up func name

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Define your message constants if they are not in a separate file
const (
	RepositoryLogger = "Repository"
	ServiceLogger    = "Service"
	ControllerLogger = "Controller"
)

// CustomFormatter adds colors and structured elements to log messages.
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format formats the log entry with colors for different parts.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Enable colors if not explicitly disabled by environment or config
	// The previous logic that forced DisableColors = false is removed.
	// logrus.TextFormatter default behavior for colors is good.

	// Define color codes
	const (
		Reset      = "\033[0m"
		Green      = "\033[32m" // Color for Info
		Yellow     = "\033[33m" // Color for Warn
		Red        = "\033[31m" // Color for Error
		Blue       = "\033[34m" // Color for Debug
		WhiteOnRed = "\033[41m" // Color for Fatal and Panic
		Cyan       = "\033[36m" // For additional info like func_name, concurrency
	)

	// Set log level color
	var logLevelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		logLevelColor = Green
	case logrus.WarnLevel:
		logLevelColor = Yellow
	case logrus.ErrorLevel:
		logLevelColor = Red
	case logrus.DebugLevel:
		logLevelColor = Blue
	case logrus.FatalLevel, logrus.PanicLevel:
		logLevelColor = WhiteOnRed
	default:
		logLevelColor = Reset // Default to no color
	}

	// Prepare additional fields for structured output
	var fieldsBuilder strings.Builder
	for key, value := range entry.Data {
		// Skip keys that are handled separately or are the main message
		if key == logrus.FieldKeyMsg || key == "level" || key == "time" {
			continue
		}
		// Special handling for 'data' field to JSON marshal it
		if key == "data" {
			dataBytes, err := json.Marshal(value)
			if err == nil {
				fieldsBuilder.WriteString(fmt.Sprintf("%s%s=%s%s ", Cyan, key, string(dataBytes), Reset))
			} else {
				fieldsBuilder.WriteString(fmt.Sprintf("%s%s=%v%s ", Cyan, key, value, Reset))
			}
		} else {
			// Add other fields with a default color (Cyan)
			fieldsBuilder.WriteString(fmt.Sprintf("%s%s=%v%s ", Cyan, key, value, Reset))
		}
	}

	concurrencyLevel := runtime.NumGoroutine()

	// Build the log message with colors for each part
	// Using entry.Message as the primary log message
	logMessage := fmt.Sprintf("%s %s[%s]%s %sconcurrency[%s%d%s]%s %s%s%s %s",
		entry.Time.Format("2006-01-02 15:04:05"),   // Timestamp
		logLevelColor, entry.Level.String(), Reset, // Level with color
		Cyan, strconv.Itoa(concurrencyLevel), Reset, // Concurrency with color
		fieldsBuilder.String(),              // Structured fields
		logLevelColor, entry.Message, Reset, // Main message with level color
	)

	return []byte(strings.TrimSpace(logMessage) + "\n"), nil
}

// InitializeLogger initializes the logger with file output and rotation
func InitializeLogger(logPath string) {
	// Ensure the logs directory exists
	logsDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory '%s': %v", logsDir, err)
	}

	// Set up log file rotation with lumberjack
	rotatingFile := &lumberjack.Logger{
		Filename:   logPath, // Log file name
		MaxSize:    10,      // Maximum size in MB before it is rotated
		MaxBackups: 5,       // Maximum number of old log files to retain
		MaxAge:     30,      // Maximum days to retain old log files
		Compress:   true,    // Compress old log files
	}

	// Set up multi-writer to log to both file and console
	multiWriter := io.MultiWriter(os.Stdout, rotatingFile)

	logrus.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			FullTimestamp: true, // Custom formatter will handle format
			ForceColors:   true, // Force colors for console, CustomFormatter will apply if enabled
		},
	})

	logrus.SetOutput(multiWriter)      // Output to both console and file
	logrus.SetLevel(logrus.DebugLevel) // Set the log level (adjust as needed)
}

// getCallerInfo gets the file and line number of the caller.
// depth: 0 for getCallerInfo itself, 1 for the function calling getCallerInfo, etc.
func getCallerInfo(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(depth + 1) // +1 because we're calling this from our wrapper
	if !ok {
		return "unknown", 0
	}
	return filepath.Base(file), line
}

// This helper is for extracting the function name
func getFunctionName(depth int) string {
	pc, _, _, ok := runtime.Caller(depth + 1) // +1 because we're calling this from our wrapper
	if !ok {
		return "unknown_func"
	}
	fullName := runtime.FuncForPC(pc).Name()
	// Trim package path from the function name
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash != -1 {
		fullName = fullName[lastSlash+1:]
	}
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot != -1 {
		return fullName[lastDot+1:] // Return just the function name
	}
	return fullName
}

// Logger logs the given data with a specified log level and service information.
// It automatically infers the calling function and file.
func Logger(data interface{}, status int, message string, logLevel logrus.Level) {
	file, line := getCallerInfo(1) // Caller of Logger
	funcName := getFunctionName(1) // Caller of Logger

	logEntry := logrus.WithFields(logrus.Fields{
		"source":    fmt.Sprintf("%s:%d", file, line),
		"func_name": funcName,
		"data":      data,
		"status":    status,
	})

	switch logLevel {
	case logrus.ErrorLevel:
		logEntry.Error(message)
	case logrus.WarnLevel:
		logEntry.Warn(message)
	case logrus.InfoLevel:
		logEntry.Info(message)
	case logrus.DebugLevel:
		logEntry.Debug(message)
	default:
		logEntry.Info(message) // Default log level will be Info
	}
}

// LoggerRequest logs details about an incoming request.
func LoggerRequest(data interface{}, action string, msg string) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logrus.WithFields(logrus.Fields{
		"source":      fmt.Sprintf("%s:%d", file, line),
		"func_name":   funcName,
		"logger_type": ControllerLogger, // Indicates this log is from a controller/request handler
		"data":        data,
		"action":      action,
		"status":      http.StatusAccepted, // Or appropriate HTTP status
	}).Info(msg) // Main message is the provided 'msg'
}

// LoggerService logs details about service layer operations.
func LoggerService(data interface{}, action string, logLevel logrus.Level) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logEntry := logrus.WithFields(logrus.Fields{
		"source":      fmt.Sprintf("%s:%d", file, line),
		"func_name":   funcName,
		"logger_type": ServiceLogger, // Indicates this log is from a service layer
		"data":        data,
		"action":      action,
		"status":      http.StatusAccepted, // Or appropriate service status
	})

	switch logLevel {
	case logrus.ErrorLevel:
		logEntry.Error(action) // 'action' can be the message for service logs
	case logrus.WarnLevel:
		logEntry.Warn(action)
	case logrus.InfoLevel:
		logEntry.Info(action)
	case logrus.DebugLevel:
		logEntry.Debug(action)
	default:
		logEntry.Info(action)
	}
}

// LoggerRepository logs details about repository layer operations.
func LoggerRepository(data interface{}, action string) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logrus.WithFields(logrus.Fields{
		"source":      fmt.Sprintf("%s:%d", file, line),
		"func_name":   funcName,
		"logger_type": RepositoryLogger, // Indicates this log is from a repository layer
		"data":        data,
		"action":      action,
		"status":      http.StatusAccepted, // Or appropriate repository status
	}).Info(action) // 'action' can be the message for repository logs
}

// ErrorLog is a convenience function for error level logs.
func ErrorLog(data interface{}, message string) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logrus.WithFields(logrus.Fields{
		"source":    fmt.Sprintf("%s:%d", file, line),
		"func_name": funcName,
		"data":      data,
		"status":    http.StatusInternalServerError,
	}).Error(message)
}

// InfoLog is a convenience function for info level logs.
func InfoLog(data interface{}, message string) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logrus.WithFields(logrus.Fields{
		"source":    fmt.Sprintf("%s:%d", file, line),
		"func_name": funcName,
		"data":      data,
		"status":    http.StatusOK,
	}).Info(message)
}

// WarnLog is a convenience function for warn level logs.
func WarnLog(data interface{}, message string) {
	file, line := getCallerInfo(1)
	funcName := getFunctionName(1)

	logrus.WithFields(logrus.Fields{
		"source":    fmt.Sprintf("%s:%d", file, line),
		"func_name": funcName,
		"data":      data,
		"status":    http.StatusAccepted,
	}).Warn(message)
}
