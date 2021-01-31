// Package log defines common structured-logging fields and functions.
package log

import (
	"fmt"
	"io"
)

// Log-entry severities.
const (
	LevelDebug   Level = Level("debug")
	LevelInfo          = Level("info")
	LevelWarn          = Level("warning")
	LevelError         = Level("error")
	LevelFatal         = Level("fatal")
	LevelPanic         = Level("panic")
	LevelInvalid       = Level("invalid")
)

// Level is a severity level.
type Level string

var stringToLevel = map[string]Level{
	"debug":   LevelDebug,
	"info":    LevelInfo,
	"warn":    LevelWarn,
	"warning": LevelWarn,
	"error":   LevelError,
	"fatal":   LevelFatal,
	"panic":   LevelPanic,
}

// Levels returns a slice of all severity levels.
func Levels() []Level {
	var levels []Level
	for _, v := range stringToLevel {
		levels = append(levels, v)
	}
	return levels
}

// ParseLevel returns a severity Level from its string value.
func ParseLevel(s string) (Level, error) {
	level, ok := stringToLevel[s]
	if !ok {
		return LevelInvalid, fmt.Errorf("invalid level string: %s", s)
	}
	return level, nil
}

// GoStdLogger are most of what's exported by Go's log package with some
// exceptions.
//
// Exceptions:
// - Flags() int
// - Output(calldepth int, s string) error
// - SetFlags(flag int)
// - SetPrefix(prefix string)
type GoStdLogger interface {
	// SetOutput sets the output destination for the registered logger.
	SetOutput(w io.Writer)

	// Print writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Print.
	Print(v ...interface{})
	// Printf writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})
	// Println writes fields+text in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Println.
	Println(v ...interface{})

	// Fatal is equivalent to Print() followed by a call to os.Exit(1).
	Fatal(v ...interface{})
	// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
	Fatalf(format string, v ...interface{})
	// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
	Fatalln(v ...interface{})

	// Panic is equivalent to Print() followed by a call to panic().
	Panic(v ...interface{})
	// Panicf is equivalent to Printf() followed by a call to panic().
	Panicf(format string, v ...interface{})
	// Panicln is equivalent to Println() followed by a call to panic().
	Panicln(v ...interface{})
}

// Fields represents field names/values as a map.
type Fields map[string]interface{}

// ExportLogFields implements FieldsExporter.
func (f Fields) ExportLogFields() Fields {
	return f
}

func (f Fields) contains(key string) bool {
	_, ok := f[key]
	return ok
}

// containsAll returns nil if all keys are present or a slice of missing keys.
func (f Fields) containsAll(keys []string) []string {
	var missing []string
	for _, v := range keys {
		if !f.contains(v) {
			missing = append(missing, v)
		}
	}
	return missing
}

// FieldsExporter provides a way for types to export loggable fields.
type FieldsExporter interface {
	ExportLogFields() Fields
}

// Logger extends GoStdLogger with levels and structured fields.
type Logger interface {
	GoStdLogger

	// SetLevel sets the minimum severity level to send to an output
	// destination.
	SetLevel(level Level) error

	// GetLevel returns the minimum severity level to send to an output
	// destination.
	GetLevel() Level

	// WithError configures a Logger to send an error field to an output
	// destination.
	WithError(err error) Logger

	// WithField configures a Logger to also send a field to an output
	// destination.
	WithField(key string, value interface{}) Logger

	// WithFields configures a Logger to also send a set of fields to an
	// output destination.
	WithFields(fields FieldsExporter) Logger

	// Debug writes fields+text in a fixed format with debug severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Print.
	Debug(args ...interface{})
	// Debugf writes fields+text in a fixed format with debug severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Printf.
	Debugf(format string, args ...interface{})
	// Debugln writes fields+text in a fixed format with debug severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Println.
	Debugln(args ...interface{})

	// Info writes fields+text in a fixed format with info severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infof writes fields+text in a fixed format with info severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Infoln writes fields+text in a fixed format with info severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Println.
	Infoln(args ...interface{})

	// Warn writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Print.
	Warn(args ...interface{})
	// Warning writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Print.
	Warning(args ...interface{})
	// Warnf writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, args ...interface{})
	// Warningf writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Printf.
	Warningf(format string, args ...interface{})
	// Warnln writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Println.
	Warnln(args ...interface{})
	// Warningln writes fields+text in a fixed format with warn severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Println.
	Warningln(args ...interface{})

	// Error writes fields+text in a fixed format with error severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorf writes fields+text in a fixed format with error severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Errorln writes fields+text in a fixed format with error severity to the registered
	// logger.  Arguments are handled in the manner of fmt.Println.
	Errorln(args ...interface{})
}
