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
	if entry.Caller != nil {
		funcName = entry.Caller.Function
	}

	// Define color codes
	const (
		Reset      = "\033[0m"
		Green      = "\033[32m"
		Yellow     = "\033[33m"
		Red        = "\033[31m"
		Blue       = "\033[34m"
		WhiteOnRed = "\033[41m"
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

	concurrencyLevel := runtime.NumGoroutine()

	logMessage := entry.Time.Format("2006-01-02 15:04:05") + " " +
		logLevelColor + "[" + entry.Level.String() + "]" + Reset + ": " +
		"[" + funcName + "]" +
		" concurrency[" + logLevelColor + strconv.Itoa(concurrencyLevel) + Reset + "]" +
		" data=" + Reset + dataJSON

	if action, ok := entry.Data["Action"].(string); ok {
		logMessage += " action=" + Blue + action + Reset
	}
	if service, ok := entry.Data["service"].(string); ok {
		logMessage += " service=" + Green + service + Reset
	}
	if status, ok := entry.Data["status"].(int); ok {
		logMessage += " status=" + logLevelColor + strconv.Itoa(status) + Reset
	}

	logMessage += " message=" + entry.Message + Reset
	return []byte(logMessage + "\n"), nil
}

// InitializeLogger initializes the logger with file output and rotation
func InitializeLogger(logPath string) {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	rotatingFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, rotatingFile)

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
			ForceColors:   true,
		},
	})
	logrus.SetOutput(multiWriter)
	logrus.SetLevel(logrus.DebugLevel)
}

// General-purpose structured logger
func Logger(data interface{}, status int, message string, level logrus.Level) {
	entry := logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": status,
	})
	switch level {
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.DebugLevel:
		entry.Debug(message)
	default:
		entry.Info(message)
	}
}

func LoggerRequest(data interface{}, action string, msg string) {
	logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": http.StatusAccepted,
		"Action": action,
	}).Info(msg)
}

func LoggerService(data interface{}, action string, level logrus.Level) {
	entry := logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": http.StatusAccepted,
		"Action": action,
	})
	switch level {
	case logrus.ErrorLevel:
		entry.Error(action)
	case logrus.WarnLevel:
		entry.Warn(action)
	case logrus.InfoLevel:
		entry.Info(action)
	case logrus.DebugLevel:
		entry.Debug(action)
	default:
		entry.Info(action)
	}
}

func LoggerRepository(data interface{}, action string) {
	logrus.WithFields(logrus.Fields{
		"logger": "Repository",
		"data":   data,
		"status": http.StatusAccepted,
		"Action": action,
	}).Info(action)
}

func ErrorLog(data interface{}, message string) {
	logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": http.StatusInternalServerError,
	}).Error(message)
}

func InfoLog(data interface{}, message string) {
	logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": http.StatusOK,
	}).Info(message)
}

func WarnLog(data interface{}, message string) {
	logrus.WithFields(logrus.Fields{
		"data":   data,
		"status": http.StatusAccepted,
	}).Warn(message)
}
