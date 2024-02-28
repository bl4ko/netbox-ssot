package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

// Default four levels for logging.
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
}

// New creates a new Logger instance, which writes to the specified destination (file) or stdout if dest is empty. It also sets the log level.
func New(dest string, logLevel int) (*Logger, error) {
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
	return &Logger{log.New(output, "", log.LstdFlags), logLevel}, nil
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

	// file prefix for logs
	filePrefix := fmt.Sprintf("%-20s", fmt.Sprintf("%s:%d", file, line))
	if l.level > DEBUG {
		filePrefix = ""
	}

	// Add your custom logging format
	logMessage := fmt.Sprintf("%s%s", filePrefix, s)

	// Print to the desired output
	l.Println(logMessage)
	return nil
}

func (l *Logger) Debug(ctx context.Context, v ...interface{}) error {
	if l.level <= DEBUG {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "DEBUG", ctx.Value(constants.CtxSourceKey), fmt.Sprint(v...)))
	}
	return nil
}

// Debugf logs a formatted debug message.
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) error {
	if l.level <= DEBUG {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "DEBUG", ctx.Value(constants.CtxSourceKey), fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) error {
	if l.level <= INFO {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "INFO", ctx.Value(constants.CtxSourceKey), fmt.Sprint(v...)))
	}
	return nil
}

// Infof logs a formatted info message.
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) error {
	if l.level <= INFO {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "INFO", ctx.Value(constants.CtxSourceKey), fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Warning(ctx context.Context, v ...interface{}) error {
	if l.level <= WARNING {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "WARNING", ctx.Value(constants.CtxSourceKey), fmt.Sprint(v...)))
	}
	return nil
}

// Warningf logs a formatted warning message.
func (l *Logger) Warningf(ctx context.Context, format string, v ...interface{}) error {
	if l.level <= WARNING {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "WARNING", ctx.Value(constants.CtxSourceKey), fmt.Sprintf(format, v...)))
	}
	return nil
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) error {
	if l.level <= ERROR {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "ERROR", ctx.Value(constants.CtxSourceKey), fmt.Sprint(v...)))
	}
	return nil
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) error {
	if l.level <= ERROR {
		return l.Output(logCallDepth, fmt.Sprintf("%-7s (%s): %s", "ERROR", ctx.Value(constants.CtxSourceKey), fmt.Sprintf(format, v...)))
	}
	return nil
}
