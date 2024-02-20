package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	DEBUG int = iota
	INFO
	WARNING
	ERROR
)

const logCallDepth = 2

type Logger struct {
	*log.Logger
	// Level of the logger (DEBUG, INFO, WARNING, ERROR).
	level int
	// Name of the logger. It will be visible in every log after log level (name)
	name string
}

// New creates a new Logger instance, which writes to the specified destination (file) or stdout if dest is empty. It also sets the log level.
func New(dest string, logLevel int, name string) (*Logger, error) {
	var output io.Writer
	if dest == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(dest)
		if err != nil {
			return nil, err
		}
		output = file
	}
	return &Logger{log.New(output, "", log.LstdFlags|log.Lshortfile), logLevel, name}, nil
}

// Custom log output function. It is used to add additional runtime information to the log message.
func (l *Logger) Output(calldepth int, s string) error {
	// Get additional runtime information
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	// Prepare the log prefix manually to include the standard log flags

	// time prefix for logs
	now := time.Now()
	timePrefix := now.Format("2006/01/02 15:04:05")

	// file prefix for logs
	filePrefix := fmt.Sprintf("%-20s", fmt.Sprintf("%s:%d", file, line))
	if l.level > DEBUG {
		filePrefix = ""
	}

	// Add your custom logging format
	logMessage := fmt.Sprintf("%s %s%s", timePrefix, filePrefix, s)

	// Print to the desired output
	l.Println(logMessage)
	return nil
}

func (l *Logger) Debug(v ...interface{}) error {
	if l.level <= DEBUG {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "DEBUG", l.name, fmt.Sprint(v...)))
	}
	return nil
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(format string, v ...interface{}) error {
	if l.level <= DEBUG {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "DEBUG", l.name, fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Info(v ...interface{}) error {
	if l.level <= INFO {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "INFO", l.name, fmt.Sprint(v...)))
	}
	return nil
}

// Infof logs a formatted info message.
func (l *Logger) Infof(format string, v ...interface{}) error {
	if l.level <= INFO {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "INFO", l.name, fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Warning(v ...interface{}) error {
	if l.level <= WARNING {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "WARNING", l.name, fmt.Sprint(v...)))
	}
	return nil
}

// Warningf logs a formatted warning message.
func (l *Logger) Warningf(format string, v ...interface{}) error {
	if l.level <= WARNING {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "WARNING", l.name, fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Error(v ...interface{}) error {
	if l.level <= ERROR {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "ERROR", l.name, fmt.Sprint(v...)))
	}
	return nil
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, v ...interface{}) error {
	if l.level <= ERROR {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "ERROR", l.name, fmt.Sprintf(format, v...)))
	}
	return nil
}
