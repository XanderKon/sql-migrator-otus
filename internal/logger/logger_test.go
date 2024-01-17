package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// simple tests.
	tests := []struct {
		level    logLevel
		expected string
	}{
		{
			logLevel{
				label: "DEBUG",
				value: 1,
			},
			"[DEBUG]",
		},
		{
			logLevel{
				label: "INFO",
				value: 2,
			},
			"[INFO]",
		},
		{
			logLevel{
				label: "WARNING",
				value: 3,
			},
			"[WARNING]",
		},
		{
			logLevel{
				label: "ERROR",
				value: 4,
			},
			"[ERROR]",
		},
	}

	for _, test := range tests {
		t.Run(test.level.label, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(test.level.label, &buf)

			logger.Log("Some test log here")

			if !strings.Contains(buf.String(), test.expected) {
				t.Errorf("Expected message has ['%s'], but got: '%s'", test.expected, buf.String())
			}
		})
	}

	// ignore weaker level
	t.Run("ignore level", func(t *testing.T) {
		var buf bytes.Buffer

		logger := New(" WaRnInG ", &buf)
		logger.Error("Some test ERROR log here")
		if !strings.Contains(buf.String(), "[ERROR]") {
			t.Errorf("Expected message has [ERROR], but got: '%s'", buf.String())
		}

		buf.Reset()
		logger.Warning("Some test WARNING log here")
		if !strings.Contains(buf.String(), "[WARNING]") {
			t.Errorf("Expected message has [WARNING], but got: '%s'", buf.String())
		}

		// should skip this level.
		buf.Reset()
		logger.Debug("Some test DEBUG log here")
		if strings.Contains(buf.String(), "[DEBUG]") {
			t.Errorf("Expected empty message, but got: '%s'", buf.String())
		}

		// should skip this level.
		buf.Reset()
		logger.Info("Some test INFO log here")
		if strings.Contains(buf.String(), "[INFO]") {
			t.Errorf("Expected empty message, but got: '%s'", buf.String())
		}
	})

	// parse template
	t.Run("parse template", func(t *testing.T) {
		var buf bytes.Buffer

		logger := New(" ErRor ", &buf)
		logger.Error("Some test ERROR %s here", "SOME INJECT INTO TEMPLATE")
		if !strings.Contains(buf.String(), "SOME INJECT INTO TEMPLATE") {
			t.Errorf("Expected message has SOME INJECT INTO TEMPLATE, but got: '%s'", buf.String())
		}
	})
}
