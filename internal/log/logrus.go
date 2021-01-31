package log

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func newLogrusLogger(out io.Writer, fields Fields, level Level) (Logger, error) {
	logrusLevel, err := logrus.ParseLevel(string(level))
	if err != nil {
		return nil, fmt.Errorf("convert level: %s", err)
	}
	logger := logrus.Logger{
		Out:       out,
		Formatter: &logrus.JSONFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrusLevel,
	}
	logrusFields := logrus.Fields{}
	for k, v := range fields {
		logrusFields[k] = v
	}
	return &logrusLogger{
		logger: &logger,
		entry:  logger.WithFields(logrusFields),
	}, nil
}

// SetLevel sets the logging severity level on the registered logger.
func (log *logrusLogger) SetLevel(level Level) error {
	logrusLevel, err := logrus.ParseLevel(string(level))
	if err != nil {
		return fmt.Errorf("convert level: %s", err)
	}
	log.logger.Level = logrusLevel
	return nil
}

// SetOutput sets the output destination for the registered logger.
func (log *logrusLogger) SetOutput(w io.Writer) {
	log.logger.Out = w
}

// GetLevel returns the minimum severity level to send to an output
// destination.
func (log *logrusLogger) GetLevel() Level {
	return Level(log.logger.Level.String())
}

// WithError
func (log *logrusLogger) WithError(err error) Logger {
	return &logrusLogger{
		logger: log.logger,
		entry:  log.entry.WithError(err),
	}
}

func (log *logrusLogger) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

// WithField
func (log *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		logger: log.logger,
		entry:  log.entry.WithField(key, value),
	}
}

// WithFields
func (log *logrusLogger) WithFields(fields FieldsExporter) Logger {
	var exported = make(logrus.Fields)
	for k, v := range fields.ExportLogFields() {
		exported[k] = v
	}
	return &logrusLogger{
		logger: log.logger,
		entry:  log.entry.WithFields(exported),
	}
}

// Print writes text+fields in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Print(v ...interface{}) {
	log.entry.Print(v...)
}

// Printf writes text+fields in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Printf(format string, v ...interface{}) {
	log.entry.Printf(format, v...)
}

// Println writes text+fields in a fixed format to the registered logger. Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Println(v ...interface{}) {
	log.entry.Println(v...)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (log *logrusLogger) Fatal(v ...interface{}) {
	log.entry.Fatal(v...)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (log *logrusLogger) Fatalf(format string, v ...interface{}) {
	log.entry.Fatalf(format, v...)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func (log *logrusLogger) Fatalln(v ...interface{}) {
	log.entry.Fatalln(v...)
}

// Panic is equivalent to Print() followed by a call to panic().
func (log *logrusLogger) Panic(v ...interface{}) {
	log.entry.Panic(v...)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func (log *logrusLogger) Panicf(format string, v ...interface{}) {
	log.entry.Panicf(format, v...)
}

// Panicln is equivalent to Println() followed by a call to panic().
func (log *logrusLogger) Panicln(v ...interface{}) {
	log.entry.Panicln(v...)
}

// Debug writes text+fields in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Debug(v ...interface{}) {
	log.entry.Debug(v...)
}

// Debugf writes text+fields in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Debugf(format string, v ...interface{}) {
	log.entry.Debugf(format, v...)
}

// Debugln writes text+fields in a fixed format with debug severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Debugln(v ...interface{}) {
	log.entry.Debugln(v...)
}

// Info writes text+fields in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Info(v ...interface{}) {
	log.entry.Info(v...)
}

// Infof writes text+fields in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Infof(format string, v ...interface{}) {
	log.entry.Infof(format, v...)
}

// Infoln writes text+fields in a fixed format with info severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Infoln(v ...interface{}) {
	log.entry.Infoln(v...)
}

// Warn writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Warn(v ...interface{}) {
	log.entry.Warn(v...)
}

// Warning writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Warning(v ...interface{}) {
	log.entry.Warning(v...)
}

// Warnf writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Warnf(format string, v ...interface{}) {
	log.entry.Warnf(format, v...)
}

// Warningf writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Warningf(format string, v ...interface{}) {
	log.entry.Warningf(format, v...)
}

// Warnln writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Warnln(v ...interface{}) {
	log.entry.Warnln(v...)
}

// Warningln writes text+fields in a fixed format with warn severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Warningln(v ...interface{}) {
	log.entry.Warningln(v...)
}

// Error writes text+fields in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Print.
func (log *logrusLogger) Error(v ...interface{}) {
	log.entry.Error(v...)
}

// Errorf writes text+fields in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Printf.
func (log *logrusLogger) Errorf(format string, v ...interface{}) {
	log.entry.Errorf(format, v...)
}

// Errorln writes text+fields in a fixed format with error severity to the registered
// logger.  Arguments are handled in the manner of fmt.Println.
func (log *logrusLogger) Errorln(v ...interface{}) {
	log.entry.Errorln(v...)
}
