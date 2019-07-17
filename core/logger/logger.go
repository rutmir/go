package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// LEVEL_DEFAULT   = 0
	levelDebug = 100
	levelInfo  = 200
	// LEVEL_NOTICE    = 300
	levelWarning = 400
	levelError   = 500
	// LEVEL_CRITICAL  = 600
	// LEVEL_ALERT     = 700
	// LEVEL_EMERGENCY = 800

	targetStackdriver = "stackdriver"
)

var (
	logDefault *log.Logger
	logFatal   *log.Logger
	logError   *log.Logger
	logWarning *log.Logger
	logInfo    *log.Logger
	logDebug   *log.Logger

	separator = "/internal"
	target    = "console"
)

// TraceItem ...
type TraceItem struct {
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// LogEntry compatible with StackDriver ...
type LogEntry struct {
	Severity int         `json:"severity"` // for StackDriver
	Message  string      `json:"message"`
	Trace    []TraceItem `json:"trace,omitempty"`
}

// JSON ...
func (le *LogEntry) JSON() string {
	b, err := json.Marshal(le)
	if err != nil {
		Err(err)
		return ""
	}
	return string(b)
}

func (le LogEntry) String() string {
	buffer := bytes.NewBufferString(le.Message)
	for _, line := range le.Trace {
		buffer.WriteString(fmt.Sprintf("\n\t-> [%s:%d] %s", line.File, line.Line, line.Name))
	}

	return buffer.String()
}

func createLogEntry(errMsg string, severity int, withTrace bool) LogEntry {
	if withTrace {
		return LogEntry{
			Severity: severity,
			Message:  errMsg,
			Trace:    extractTrace(),
		}
	}

	return LogEntry{
		Severity: severity,
		Message:  errMsg,
	}
}

func formatOutput(str string, file string, line int) string {
	return fmt.Sprintf("[%s:%s] - %s", file, strconv.Itoa(line), str)
}

func extractString(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	if first, ok := v[0].(string); ok {
		if len(v) > 1 {
			return fmt.Sprintf(first, v[1:]...)
		}
		return first
	}

	return fmt.Sprint(v...)
}

func getCallers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[0:n]
}

func extractTrace() []TraceItem {
	result := make([]TraceItem, 0)
	callers := getCallers(3)

	for idx, item := range callers {
		if idx < 3 {
			continue
		}

		fn := runtime.FuncForPC(item)
		if fn == nil {
			continue
		}

		file, line := fn.FileLine(item)
		splitFile := strings.Split(file, separator)
		if len(splitFile) < 2 {
			continue
		}

		name := "unknown"
		splitName := strings.Split(fn.Name(), ".")
		if len(splitName) > 0 {
			name = splitName[len(splitName)-1]
		}

		result = append(result, TraceItem{Name: name, File: splitFile[len(splitFile)-1], Line: line})
	}

	return result
}

func extractPath() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = ""
		line = 0
	} else {
		split := strings.Split(file, separator)
		if len(split) > 0 {
			file = split[len(split)-1]
		}
	}
	return file, line
}

func defaultLogEntry(msg string, level int, withTrace bool, logger *log.Logger) {
	le := createLogEntry(msg, level, withTrace)

	switch target {
	case targetStackdriver:
		logDefault.Println(le.JSON())
	default:
		logger.Println(le)
	}
}

// HTTPLogger is logger for every http requests
func HTTPLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf("%s\t%s\t%s\t%s\n", r.Method, r.RequestURI, name, time.Since(start))
	})
}

// Fatal handler for LevelFatal
func Fatal(v ...interface{}) {
	str := extractString(v...)
	file, line := extractPath()

	logFatal.Println(formatOutput(str, file, line))
	os.Exit(1)
}

// Err handler for LevelError
func Err(v ...interface{}) {
	str := extractString(v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelError, false, logError)
}

// Errf ...
func Errf(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelError, false, logError)
}

// ErrTrace display error with trace
func ErrTrace(err error) {
	defaultLogEntry(fmt.Sprintf("%v", err), levelError, true, logError)
}

// ErrTracef display error with trace
func ErrTracef(format string, v ...interface{}) {
	defaultLogEntry(fmt.Sprintf(format, v...), levelError, true, logError)
}

// Warn ...
func Warn(v ...interface{}) {
	str := extractString(v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelWarning, false, logWarning)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelWarning, false, logWarning)
}

// WarnTrace ...
func WarnTrace(err error) {
	defaultLogEntry(fmt.Sprintf("%v", err), levelWarning, true, logWarning)
}

// WarnTracef ...
func WarnTracef(format string, v ...interface{}) {
	defaultLogEntry(fmt.Sprintf(format, v...), levelWarning, true, logWarning)
}

// Info ...
func Info(v ...interface{}) {
	str := extractString(v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelInfo, false, logInfo)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	str := fmt.Sprintf(format, v...)
	file, line := extractPath()

	defaultLogEntry(formatOutput(str, file, line), levelInfo, false, logInfo)
}

// InfoTrace ...
func InfoTrace(err error) {
	defaultLogEntry(fmt.Sprintf("%v", err), levelInfo, true, logInfo)
}

// InfoTracef ...
func InfoTracef(format string, v ...interface{}) {
	defaultLogEntry(fmt.Sprintf(format, v...), levelInfo, true, logInfo)
}

// Debug ...
func Debug(v ...interface{}) {
	defaultLogEntry(extractString(v...), levelDebug, false, logDebug)
}

// Debugf ...
func Debugf(format string, v ...interface{}) {
	defaultLogEntry(fmt.Sprintf(format, v...), levelDebug, false, logDebug)
}

// Initialize ...
func Initialize(pathSeparator, logTarget string) {
	if len(pathSeparator) > 0 {
		separator = pathSeparator
	}

	if len(logTarget) > 0 {
		target = logTarget
	}

	logFatal = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime)
	logError = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	logWarning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	logInfo = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logDebug = log.New(os.Stdout, "DEBUG: ", 0)
	logDefault = log.New(os.Stdout, "", 0)
}
