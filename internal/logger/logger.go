package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	DEBUG int = iota
	INFO
	WARNING
	ERROR
)

type Logger struct {
	*log.Logger
	logLevel int
}

// New creates a new Logger instance, which writes to the specified destination (file) or stdout if dest is empty. It also sets the log level.
func New(dest string, logLevel int) (*Logger, error) {
	var output io.Writer
	if dest == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(fmt.Sprintf("/var/log/%s.log", dest))
		if err != nil {
			return nil, err
		}
		output = file
	}
	return &Logger{log.New(output, "", log.LstdFlags), logLevel}, nil
}

func (l *Logger) Debug(v ...interface{}) error {
	if l.logLevel <= DEBUG {
		err := l.Output(2, fmt.Sprintf("DEBUG: %s", fmt.Sprint(v...)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Info(v ...interface{}) error {
	if l.logLevel <= INFO {
		err := l.Output(2, fmt.Sprintf("INFO: %s", fmt.Sprint(v...)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Warning(v ...interface{}) error {
	if l.logLevel <= WARNING {
		err := l.Output(2, fmt.Sprintf("WARNING: %s", fmt.Sprint(v...)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logger) Error(v ...interface{}) error {
	if l.logLevel <= ERROR {
		err := l.Output(2, fmt.Sprintf("ERROR: %s", fmt.Sprint(v...)))
		if err != nil {
			return err
		}
	}
	return nil
}
