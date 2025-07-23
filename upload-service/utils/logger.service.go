package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// CustomFormatter adds colors to log levels and other parts of the log message.
type CustomFormatter struct {
	logrus.TextFormatter
}

// Format formats the log entry with colors for different parts.
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	funcName := ""

	pc, _, _, ok := runtime.Caller(6) // Adjust the caller level as needed
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	// Enable colors if the terminal supports it
	if f.TextFormatter.DisableColors || os.Getenv("TERM") == "" {
		f.TextFormatter.DisableColors = false
	} else {
		f.TextFormatter.DisableColors = false
	}

	// Define color codes
	const (
		Reset      = "\033[0m"
		Green      = "\033[32m" // Color for Info
		Yellow     = "\033[33m" // Color for Warn
		Red        = "\033[31m" // Color for Error
		Blue       = "\033[34m" // Color for Debug
		WhiteOnRed = "\033[41m" // Color for Fatal and Panic
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
		logLevelColor = Reset
	}

	// Prepare fields with colors
	var dataJSON string
	if entry.Data != nil {
		dataJSONBytes, _ := json.Marshal(entry.Data)
		dataJSON = string(dataJSONBytes)
	}

	// Prepare service and status safely
	service := entry.Data["service"]
	status := entry.Data["status"]

	concurrencyLevel := runtime.NumGoroutine()
	// Convert service and status to string, checking if they are set
	serviceStr, serviceOk := service.(string)
	statusStr, statusOk := status.(int)
	logColor := ""
	switch entry.Message {
	case string(SuccessMessage):
		logColor = Green
	case string(NotFoundMessage):
		logColor = Yellow
	case string(ErrorMessage):
		logColor = Red
	case string(ExecuteMessage):
		logColor = Green
	case string(RequestMessage):
		logColor = Green
	case string(ResponseMessage):
		logColor = Green
	default:
		logColor = Red
	}
	// Build the log message with colors for each part
	logMessage := entry.Time.Format("2006-01-02 15:04:05") + " " +
		logLevelColor + "[" + entry.Level.String() + "]" + Reset + ": " +
		logColor + "[" + funcName + "]" + Reset +
		" concurrency[" + logColor + strconv.Itoa(concurrencyLevel) + Reset + "]" +
		" data=" + Reset + dataJSON + Reset +
		" action=" + logColor + entry.Message + Reset +
		" message=" + logColor + entry.Message + Reset

	if serviceOk {
		logMessage += " service=" + Green + serviceStr + Reset
	}

	if statusOk {
		logMessage += " status=" + logColor + strconv.Itoa(statusStr) + Reset // Correctly convert int to string
	}

	// Return the log message as bytes
	return []byte(logMessage + "\n"), nil
}

// InitializeLogger initializes the logger with file output and rotation
func InitializeLogger(logPath string) {
	logFile := logPath // Set the log file path relative to the project root

	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Set up log file rotation with lumberjack
	rotatingFile := &lumberjack.Logger{
		Filename:   logFile, // Log file name
		MaxSize:    10,      // Maximum size in MB before it is rotated
		MaxBackups: 5,       // Maximum number of old log files to retain
		MaxAge:     30,      // Maximum days to retain old log files
		Compress:   true,    // Compress old log files
	}

	// Set up multi-writer to log to both file and console
	multiWriter := io.MultiWriter(os.Stdout, rotatingFile)

	logrus.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
			ForceColors:   true, // Forces colors if supported
		},
	})

	logrus.SetOutput(multiWriter)      // Output to both console and file
	logrus.SetLevel(logrus.DebugLevel) // Set the log level (adjust as needed)
}

// Logger logs the given data with a specified log level and service information.
func Logger(data interface{}, status int, message string, logLevel logrus.Level) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  status,
		"message": message,
	})

	// Log based on the specified log level
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

func LoggerRequest(data interface{}, action string, msg string) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  http.StatusAccepted,
		"Action":  action,
		"message": msg,
	})
	logEntry.Info(msg)
}

func LoggerService(data interface{}, action string, logLevel logrus.Level) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  http.StatusAccepted,
		"Action":  action,
		"message": "",
	})
	// Log based on the specified log level
	switch logLevel {
	case logrus.ErrorLevel:
		logEntry.Error(action)
	case logrus.WarnLevel:
		logEntry.Warn(action)
	case logrus.InfoLevel:
		logEntry.Info(action)
	case logrus.DebugLevel:
		logEntry.Debug(action)
	default:
		logEntry.Info(action) // Default log level will be Info
	}
}

func LoggerRepository(data interface{}, action string) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"logger":  "Repository",
		"data":    data,
		"status":  http.StatusAccepted,
		"Action":  action,
		"message": "",
	})
	logEntry.Info(action)
}

func ErrorLog(data interface{}, message string) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  http.StatusInternalServerError,
		"message": message,
	})
	logEntry.Error(message)
}

func InfoLog(data interface{}, message string) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  http.StatusOK,
		"message": message,
	})
	logEntry.Info(message)
}

func WarnLog(data interface{}, message string) {
	// Create a structured log entry with additional fields
	logEntry := logrus.WithFields(logrus.Fields{
		"data":    data,
		"status":  http.StatusAccepted,
		"message": message,
	})
	logEntry.Warn(message)
}
