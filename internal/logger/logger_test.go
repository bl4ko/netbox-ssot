package logger

import (
	"testing"
)

func TestLoggerCreationForStdout(t *testing.T) {
	// We create logger with dest: "" (stdout) and logLevel: DEBUG
	logger, err := New("", DEBUG, "test")
	if err != nil {
		t.Errorf("Error creating logger: %v", err)
	}

	if logger == nil {
		t.Errorf("Logger is nil")
	}

	err = logger.Info("Test INFO")
	if err != nil {
		t.Errorf("Error writing to logger on level INFO: %v", err)
	}

	testString := "info"
	err = logger.Infof("Test %s", testString)
	if err != nil {
		t.Errorf("Error writing to logger on level INFO: %v", err)
	}

	err = logger.Debug("Test DEBUG")
	if err != nil {
		t.Errorf("Error writing to logger on level DEBUG: %v", err)
	}

	err = logger.Debugf("Test %s", testString)
	if err != nil {
		t.Errorf("Error writing to logger on level DEBUG: %v", err)
	}

	err = logger.Warning("Test WARNING")
	if err != nil {
		t.Errorf("Error writing to logger on level WARNING: %v", err)
	}
	err = logger.Warningf("Test %s", testString)
	if err != nil {
		t.Errorf("Error writing to logger on level WARNING: %v", err)
	}

	err = logger.Error("Test ERROR")
	if err != nil {
		t.Errorf("Error writing to logger on level ERROR: %v", err)
	}
	err = logger.Errorf("Test %s", testString)
	if err != nil {
		t.Errorf("Error writing to logger on level ERROR: %v", err)
	}
}

// TODO: Test for custom file logging
