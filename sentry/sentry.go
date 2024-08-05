package sentry

import (
	"encoding/json"
	"github.com/catalystcommunity/app-utils-go/env"
	"github.com/catalystcommunity/app-utils-go/errorutils"
	"github.com/catalystcommunity/app-utils-go/logging"
	"github.com/getsentry/sentry-go"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

var SentryEnabled = env.GetEnvAsBoolOrDefault("SENTRY_ENABLED", "false")
var SentryDebug = env.GetEnvAsBoolOrDefault("SENTRY_DEBUG", "false")
var SampleRate = env.GetEnvAsFloatOrDefault("SENTRY_SAMPLE_RATE", "1.0")
var TracesSampleRate = env.GetEnvAsFloatOrDefault("SENTRY_TRACES_SAMPLE_RATE", "1.0")
var SentryDsn = env.GetEnvOrDefault("SENTRY_DSN", "")
var AdditionalSentryTags = env.GetEnvOrDefault("ADDITIONAL_SENTRY_TAGS", "{}")

func MaybeInitSentry(options sentry.ClientOptions, hook *SentryLogrusHook) {
	if SentryEnabled {
		options.SampleRate = SampleRate
		options.TracesSampleRate = TracesSampleRate
		options.Debug = SentryDebug
		// I'm setting AttachStacktrace here because I can't imagine a time when we wouldn't want this and it will simplify implementation code
		// if we ever want this to be optional, we can change it.
		options.AttachStacktrace = true
		err := sentry.Init(options)
		if err != nil {
			// init doesn't actually care about connections. Errors from init are things like improperly formatted DSN,
			// fat fingered settings, stuff like that. So I'm intentionally dying here if sentry is enabled but misconfigured.
			logging.Log.WithError(err).Fatal("error initializing sentry")
		}
		addSentryLogrusHook(hook)
	}
}

func addSentryLogrusHook(hook *SentryLogrusHook) {
	// get additional tags from env var
	tags := map[string]string{}
	err := json.Unmarshal([]byte(AdditionalSentryTags), &tags)
	errorutils.PanicOnErr(nil, "ADDITIONAL_SENTRY_HOOKS is invalid, it should be json that maps to a map[string]string", err)
	// initialize hook
	if hook == nil {
		hook = NewDefaultSentryLogrusHook()
	}
	// defaults, these are set in every pod so we get them for free
	tags["namespace"] = env.GetEnvOrDefault("POD_NAMESPACE", "local")
	tags["pod_name"] = env.GetEnvOrDefault("HOSTNAME", "local")
	// add additional tags to the hook
	if hook.tags == nil {
		hook.tags = map[string]string{}
	}
	for key, value := range tags {
		hook.tags[key] = value
	}
	// add hook to logrus
	logging.Log.Hooks.Add(hook)
}

func AutoRecoverWithCapture(log *logrus.Entry, message string) error {
	recoverErr := recover()
	return RecoverWithCapture(log, message, recoverErr)
}

func RecoverWithCapture(log *logrus.Entry, message string, recover interface{}) error {
	err := errorutils.RecoverErr(recover)
	errorutils.LogOnErr(log, message, err)
	return err
}

type SentryLogrusHook struct {
	levels []logrus.Level
	tags   map[string]string
}

func NewDefaultSentryLogrusHook() *SentryLogrusHook {
	return &SentryLogrusHook{
		levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		},
		tags: map[string]string{},
	}
}
func NewSentryLogrusHook(levels []logrus.Level) *SentryLogrusHook {
	return &SentryLogrusHook{levels: levels}
}

func (s *SentryLogrusHook) Levels() []logrus.Level {
	return s.levels
}

func (s *SentryLogrusHook) Fire(entry *logrus.Entry) error {
	var err error
	logEntryErr := entry.Data["error"]
	if logEntryErr == nil {
		err = errorx.InternalError.New(entry.Message)
	} else {
		err = logEntryErr.(error)
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTags(s.tags)
		sentry.CaptureException(err)
	})
	return nil
}
