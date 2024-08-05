package logging

import (
	"strings"

	"github.com/catalystcommunity/app-utils-go/env"
	"github.com/sirupsen/logrus"
)

var LogLevel = env.GetEnvOrDefault("LOG_LEVEL", "INFO")
var LogFmt = env.GetEnvAsBoolOrDefault("LOG_FMT", "false")
var Log = logrus.New()

func init() {
	// logrus
	if LogFmt {
		// default uses `logfmt` if terminal is not tty. This is nice for machines but bad for human readability
		Log.Formatter = &logrus.TextFormatter{}
	} else {
		// force colors is more human readable
		Log.Formatter = &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		}
	}
	// levels
	setLogLevels()
}

func setLogLevels() {
	var logrusLogLevel logrus.Level
	switch strings.ToLower(LogLevel) {
	case "panic":
		logrusLogLevel = logrus.PanicLevel
	case "fatal":
		logrusLogLevel = logrus.FatalLevel
	case "error":
		logrusLogLevel = logrus.ErrorLevel
	case "warn":
		logrusLogLevel = logrus.WarnLevel
	case "info":
		logrusLogLevel = logrus.InfoLevel
	case "debug":
		logrusLogLevel = logrus.DebugLevel
	case "trace":
		logrusLogLevel = logrus.TraceLevel
	default:
		logrusLogLevel = logrus.InfoLevel
	}
	Log.Level = logrusLogLevel
}
