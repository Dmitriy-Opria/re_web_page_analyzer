package log

import (
	"fmt"
	"io"
	"os"
)

// DefaultMinLevel is the default minimum logging level.
const DefaultMinLevel = LevelInfo

var (
	// defaultLogger is currentLogger's initial value and can never
	// set to a nil value.
	defaultLogger Logger

	// currentLogger is called by all exported logging functions
	// and can never be set to a nil value.
	currentLogger Logger

	// output logs to this stream.
	output io.Writer
)

// init configures a default logger at LevelInfo and an initial output
// destination set to os.Stderr.
func init() {
	output = os.Stdout
	var err error
	defaultLogger, err = newLogrusLogger(output, Fields{}, DefaultMinLevel)
	if err != nil {
		panic(fmt.Sprintf("log.init() create default logger: %s", err))
	}
	currentLogger = defaultLogger

	// post invariant checks
	if defaultLogger == nil {
		panic("no defaultLogger set!")
	}
	if currentLogger == nil {
		panic("no currentLogger set!")
	}
}

// GetDefaultLogger returns the default logger setup during package init.
func GetDefaultLogger() Logger {
	return defaultLogger
}

// GetLogger returns the configured standard logger.
func GetLogger() Logger {
	return currentLogger
}

// SetLogger sets the registered logger.  There must always be a
// valid registered logger so the log parameter is ignored if it is nil.
func SetLogger(log Logger) {
	if log != nil {
		currentLogger = log
	}
}

// NewCommandLogger returns a logger configured for CLI programs.  The logger's
// output is set to the current output destination by calling GetOutput().
func NewCommandLogger(minimumLevel Level) (Logger, error) {
	return newLogrusLogger(GetOutput(), Fields{}, minimumLevel)
}

var requiredLoggerFields = []string{
	FieldCategory,
	FieldService,
}

// NewServiceLogger returns a logger configured with a minimum set of standard
// fields.  The logger's output is set to the current output destination by
// calling GetOutput().  All service loggers require these structured fields:
// 	FieldCategory
//  	FieldService
func NewServiceLogger(fields Fields, minimumLevel Level) (Logger, error) {
	if missing := fields.containsAll(requiredLoggerFields); len(missing) > 0 {
		return nil, fmt.Errorf("missing required fields: %s", missing)
	}
	return newLogrusLogger(GetOutput(), fields, minimumLevel)
}

// SetLevel sets the logging severity level on the registered logger.
// Both the default and current logger are reconfigured to use the same
// level.
func SetLevel(level Level) {
	if level != LevelInvalid {
		_ = currentLogger.SetLevel(level)
		_ = defaultLogger.SetLevel(level)
	}
}

// GetLevel returns the minimum severity level to send to an output
// destination.
func GetLevel() Level {
	return currentLogger.GetLevel()
}

// SetOutput sets the output destination for the registered logger.
// Both the default and current logger are reconfigured to use the same
// output destination.
func SetOutput(w io.Writer) {
	output = w
	currentLogger.SetOutput(output)
	defaultLogger.SetOutput(output)
}

// GetOutput gets the output destination of the registered logger.
// Defaults to os.Stderr.
func GetOutput() io.Writer {
	return output
}

// WithError returns a Logger configured to write an error field.
func WithError(err error) Logger {
	return currentLogger.WithError(err)
}

// WithField returns a Logger configured to write the given field.
func WithField(key string, value interface{}) Logger {
	return currentLogger.WithField(key, value)
}

// WithFields returns a Logger configured to write the given fields.
func WithFields(fields FieldsExporter) Logger {
	return currentLogger.WithFields(fields)
}

// Print writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	currentLogger.Print(v...)
}

// Printf writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	currentLogger.Printf(format, v...)
}

// Println writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	currentLogger.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	currentLogger.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	currentLogger.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	currentLogger.Fatalln(v...)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	currentLogger.Panic(v...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	currentLogger.Panicf(format, v...)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	currentLogger.Panicln(v...)
}

// Debug writes fields+text in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	currentLogger.Debug(v...)
}

// Debugf writes fields+text in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	currentLogger.Debugf(format, v...)
}

// Debugln writes fields+text in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	currentLogger.Debugln(v...)
}

// Info writes fields+text in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	currentLogger.Info(v...)
}

// Infof writes fields+text in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	currentLogger.Infof(format, v...)
}

// Infoln writes fields+text in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	currentLogger.Infoln(v...)
}

// Warn writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func Warn(v ...interface{}) {
	currentLogger.Warn(v...)
}

// Warning writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	currentLogger.Warning(v...)
}

// Warnf writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, v ...interface{}) {
	currentLogger.Warnf(format, v...)
}

// Warningf writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	currentLogger.Warningf(format, v...)
}

// Warnln writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func Warnln(v ...interface{}) {
	currentLogger.Warnln(v...)
}

// Warningln writes fields+text in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	currentLogger.Warningln(v...)
}

// Error writes fields+text in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	currentLogger.Error(v...)
}

// Errorf writes fields+text in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	currentLogger.Errorf(format, v...)
}

// Errorln writes fields+text in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	currentLogger.Errorln(v...)
}
