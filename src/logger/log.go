package log

// This is a dropin replacement for 	"google.golang.org/appengine/log"

import (
	"context"
	"log"
)

type logLevel string

const (
	debug    logLevel = "DEBUG"
	info              = "INFO"
	warning           = "WARNING"
	error             = "ERROR"
	critical          = "CRITICAL"
)

// Debugf formats its arguments according to the format, analogous to fmt.Printf,
// and records the text as a log message at Debug level. The message will be associated
// with the request linked with the provided context.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logf(debug, format, args...)
}

// Infof is like Debugf, but at Info level.
func Infof(ctx context.Context, format string, args ...interface{}) {
	logf(info, format, args...)
}

// Warningf is like Debugf, but at Warning level.
func Warningf(ctx context.Context, format string, args ...interface{}) {
	logf(warning, format, args...)
}

// Errorf is like Debugf, but at Error level.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logf(error, format, args...)
}

// Criticalf is like Debugf, but at Critical level.
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	logf(critical, format, args...)
}

func logf(level logLevel, format string, args ...interface{}) {
	log.Printf(string(level)+": "+format, args...)
}
