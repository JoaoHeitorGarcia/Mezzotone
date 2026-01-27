package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileLogger struct {
	path string
	file *os.File
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

	logger = &FileLogger{path: abs, file: f}
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

func (l *FileLogger) WriteLine(level string, msg string) error {
	if l == nil || l.file == nil {
		return fmt.Errorf("logger not initialized (call services.InitLogger first)")
	}

	ts := time.Now().Format(time.RFC3339)
	_, err := fmt.Fprintf(l.file, "%s [%s] %s\n", ts, level, msg)
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
