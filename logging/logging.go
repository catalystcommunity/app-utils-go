package logging

import (
	"github.com/catalystsquad/app-utils-go/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
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
	// zerolog
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// levels
	setLogLevels()
}

func setLogLevels() {
	var zeroLogLevel zerolog.Level
	var logrusLogLevel logrus.Level
	switch strings.ToLower(LogLevel) {
	case "panic":
		zeroLogLevel = zerolog.PanicLevel
		logrusLogLevel = logrus.PanicLevel
	case "fatal":
		zeroLogLevel = zerolog.FatalLevel
		logrusLogLevel = logrus.FatalLevel
	case "error":
		zeroLogLevel = zerolog.ErrorLevel
		logrusLogLevel = logrus.ErrorLevel
	case "warn":
		zeroLogLevel = zerolog.WarnLevel
		logrusLogLevel = logrus.WarnLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
		logrusLogLevel = logrus.InfoLevel
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
		logrusLogLevel = logrus.DebugLevel
	case "trace":
		zeroLogLevel = zerolog.TraceLevel
		logrusLogLevel = logrus.TraceLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
		logrusLogLevel = logrus.InfoLevel
	}
	zerolog.SetGlobalLevel(zeroLogLevel)
	Log.Level = logrusLogLevel
}
