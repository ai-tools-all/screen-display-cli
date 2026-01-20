package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func New(verbose bool) (*logrus.Logger, error) {
	logger := logrus.New()

	level := logrus.InfoLevel
	if verbose {
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)

	stdoutFormatter := &logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    false,
		ForceColors:      true,
	}

	logDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "dmon")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(logDir, "dmon.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	fileFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(stdoutFormatter)

	logger.AddHook(&dualOutputHook{
		file:          file,
		fileFormatter: fileFormatter,
		fileLevel:     logrus.DebugLevel,
	})

	return logger, nil
}

type dualOutputHook struct {
	file          io.Writer
	fileFormatter logrus.Formatter
	fileLevel     logrus.Level
}

func (h *dualOutputHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *dualOutputHook) Fire(entry *logrus.Entry) error {
	if entry.Level > h.fileLevel {
		return nil
	}

	line, err := h.fileFormatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.file.Write(line)
	return err
}
