package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileLogger struct {
	path     string
	file     *os.File
	maxBytes int64
}

var logger *FileLogger

func InitLogger(path string) error {
	if logger != nil {
		return nil
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(abs, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	logger = &FileLogger{path: abs, file: f, maxBytes: int64(10 * 1_000 * 1_000)}
	return nil
}

func Logger() *FileLogger {
	return logger
}

func (l *FileLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *FileLogger) WriteLine(level, msg string) error {
	if l == nil || l.file == nil {
		return fmt.Errorf("logger not initialized")
	}

	if err := l.maybeTruncate(); err != nil {
		return err
	}

	ts := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(l.file, "%s [%s] %s\n", ts, level, msg)
	return err
}

func (l *FileLogger) maybeTruncate() error {
	if l.maxBytes <= 0 {
		return nil
	}

	info, err := l.file.Stat()
	if err != nil {
		return err
	}
	if info.Size() <= l.maxBytes {
		return nil
	}

	if err := l.file.Close(); err != nil {
		return err
	}

	f, err := os.OpenFile(l.path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	l.file, err = os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	return err
}

func (l *FileLogger) Info(msg string) error  { return l.WriteLine("INFO", msg) }
func (l *FileLogger) Warn(msg string) error  { return l.WriteLine("WARN", msg) }
func (l *FileLogger) Error(msg string) error { return l.WriteLine("ERROR", msg) }

func (l *FileLogger) Infof(format string, args ...any) error {
	return l.Info(fmt.Sprintf(format, args...))
}
func (l *FileLogger) Errorf(format string, args ...any) error {
	return l.Error(fmt.Sprintf(format, args...))
}
