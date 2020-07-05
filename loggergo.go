package loggergo

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// WithJSONFormatter set the formatter of the given logger to the JSON
// formatter with the following field keys:
//
// * time -> @timestamp
// * message -> message
// * level -> level_name
func WithJSONFormatter(l *logrus.Logger) *logrus.Logger {
	l.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "level_name",
		},
	})

	return l
}

// WithLevel sets a string log level to the given logger. If there is an error
// parsing the log level, it sets the level to Warn, then it logs and returns.
func WithLevel(l *logrus.Logger, level string) *logrus.Logger {
	previous := l.GetLevel()
	l.SetLevel(logrus.WarnLevel)
	newLevel, err := logrus.ParseLevel(level)
	if err != nil {
		l.
			WithFields(logrus.Fields{"previous": previous.String(), "desired": level}).
			WithError(err).
			Warn("could not change logger level, setting level to warning")

		return l
	}

	l.SetLevel(newLevel)

	return l
}

// WithTrace attaches trace information to the logger if it's available
func WithTrace(l logrus.FieldLogger, traceID, spanID string) logrus.FieldLogger {
	l = l.
		WithField("traceID", traceID).
		WithField("spanID", spanID)

	uid, err := toUUIDString(traceID)
	if err != nil {
		l.
			WithError(err).
			Warn("could not format trace ID to UUID v4")

		return l
	}

	l = l.WithField("traceUUID", uid)

	return l
}

func toUUIDString(id string) (string, error) {
	uid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		return id, err
	}

	return uid.String(), nil
}
