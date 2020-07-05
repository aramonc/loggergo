package loggergo_test

import (
	"math/rand"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/aramonc/loggergo"
)

func TestWithJSONFormatterSetsTheJSONFormatter(t *testing.T) {
	l := logrus.New()
	ln := loggergo.WithJSONFormatter(l)

	assert.IsType(t, &logrus.JSONFormatter{}, ln.Formatter)
}

func TestWithLevelSetsCorrectLevel(t *testing.T) {
	l := logrus.New()
	ln := loggergo.WithLevel(l, "debug")

	assert.True(t, ln.IsLevelEnabled(logrus.DebugLevel))
}

func TestWithLevelRevertsToPreviousLevelOnError(t *testing.T) {
	l, h := test.NewNullLogger()

	l.SetLevel(logrus.FatalLevel)

	ln := loggergo.WithLevel(l, "dbug")

	assert.False(t, ln.IsLevelEnabled(logrus.DebugLevel))
	assert.True(t, ln.IsLevelEnabled(logrus.WarnLevel))
	assert.Equal(t, logrus.WarnLevel, h.LastEntry().Level)
}

func TestWithTraceSetsIDFields(t *testing.T) {
	tID := make([]byte, 16)
	mtID := make([]byte, 12)
	sID := make([]byte, 8)

	rand.Read(tID)
	rand.Read(mtID)
	rand.Read(sID)

	type (
		argsType struct {
			trace []byte
			span  []byte
		}

		expectedType struct {
			logRecords int
			numFields  int
		}
	)

	scenarios := []struct {
		description string
		args        argsType
		expected    expectedType
	}{
		{
			description: "succeeds and sets 3 fields to the entry",
			args: argsType{
				trace: tID,
				span:  sID,
			},
			expected: expectedType{
				logRecords: 0,
				numFields:  3,
			},
		},
		{
			description: "fails to set UUID field and logs the failure with a warning",
			args: argsType{
				trace: mtID,
				span:  sID,
			},
			expected: expectedType{
				logRecords: 1,
				numFields:  2,
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.description, func(t *testing.T) {
			l, h := test.NewNullLogger()
			ln := loggergo.WithTrace(l, string(s.args.trace), string(s.args.span))

			entry := ln.(*logrus.Entry)

			assert.Len(t, h.Entries, s.expected.logRecords)
			assert.Len(t, entry.Data, s.expected.numFields)
		})
	}

}
