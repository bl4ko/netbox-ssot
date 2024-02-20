package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerCreationForCustomFile(t *testing.T) {
	// We create logger with dest: "" (stdout) and logLevel: DEBUG
	logger, err := New("/tmp/test", DEBUG, "test")
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

func TestLoggerCreationWithWrongPath(t *testing.T) {
	// We create logger with dest: "" (stdout) and logLevel: DEBUG
	wrongFileName := "///wrongPath///"
	_, err := New(wrongFileName, DEBUG, "test")
	if err == nil {
		t.Errorf("Error creating logger: %v", err)
	}
}

func TestHighLoggerLevel(t *testing.T) {
	l, err := New("", 1, "test")
	if err != nil {
		t.Errorf("Error creating logger: %v", err)
	}
	buffer := new(bytes.Buffer)
	l.Logger.SetOutput(buffer)
	l.Debug("Test DEBUG")
	l.Debugf("Test DEBUG")
	// we need to ensure that the buffer is empty (no output)
	if buffer.String() != "" {
		t.Errorf("Buffer should be empty")
	}
	l.level = 2
	l.Info("Test INFO")
	l.Infof("Test INFO")
	if buffer.String() != "" {
		t.Errorf("Buffer should be empty")
	}
	l.level = 3
	l.Warning("Test WARNING")
	l.Warningf("Test WARNING")
	if buffer.String() != "" {
		t.Errorf("Buffer should be empty")
	}
	l.level = 4
	l.Error("Test ERROR")
	l.Errorf("Test ERROR")
	if buffer.String() != "" {
		t.Errorf("Buffer should be empty")
	}
}

func TestLogPrefixForDebugOnly(t *testing.T) {
	l, err := New("", 1, "test")
	if err != nil {
		t.Errorf("Error creating logger: %v", err)
	}

	// Create a buffer to hold stdout
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	l.Output(2, "test message")

	output := buf.String()

	// Check if the output contains the file prefix
	if strings.Contains(output, "logger_test.go:") {
		t.Errorf("Expected no file prefix, but got one")
	}
}

func TestOutputUnknownCaller(t *testing.T) {
	l, err := New("", 0, "test")
	if err != nil {
		t.Errorf("Error creating logger: %v", err)
	}

	// Create a buffer to hold stdout
	buf := new(bytes.Buffer)
	l.SetOutput(buf)

	// Larger call depth than stack
	largCallDepth := 999999
	l.Output(largCallDepth, "test message")

	output := buf.String()

	// Check if the output contains "???"
	if !strings.Contains(output, "???") {
		t.Errorf("Expected \"???\" in output, but didn't get it")
	}
}
