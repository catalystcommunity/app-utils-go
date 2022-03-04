package errorutils

import (
	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

func RecoverErr(recover interface{}) error {
	if recover == nil {
		return nil
	}
	// switch handles things like panic("oh no")
	switch recover.(type) {
	case string:
		return errorx.InternalError.New(recover.(string))
	default:
		return recover.(error)
	}
}

func LogOnErr(logEntry *logrus.Entry, message string, err error) {
	if err != nil {
		if logEntry == nil {
			logEntry = logging.Log.WithFields(nil)
		}
		err = errorx.Decorate(err, message)
		logEntry.WithError(err).Errorf("Error: %+v", err)
	}
}

//PanicWithTrace will log a formatted stack trace at the panic level, which will capture the error if sentry is enabled
func PanicOnErr(logEntry *logrus.Entry, message string, err error) {
	LogOnErr(logEntry, message, err)
	if err != nil {
		panic(err)
	}
}
