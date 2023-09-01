package logger

import (
	"github.com/sirupsen/logrus"
)

var Default Logger

func New() *logrus.Logger {
	return logrus.New()
}

func init() {
	log := New()
	// Default.SetOutput(io.MultiWriter(Default.Out, output))
	// Default.SetOutput(output)
	// Default.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors: true,
	// })
	// Default.SetFormatter(&Formatter{
	// 	// HideKeys:        true,
	// 	TimestampFormat: "2006-01-02 15:04:05.000",
	// 	NoColors:        true,
	// })
	log.SetLevel(logrus.DebugLevel)
	Default = log
}

func SetLevel(lvstr string) {
	lv, err := logrus.ParseLevel(lvstr)
	if err != nil {
		Default.Error(err)
	} else {
		Default.(*logrus.Logger).SetLevel(lv)
	}
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	Default.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	Default.Info(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	Default.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	Default.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	Default.Error(args...)
}
