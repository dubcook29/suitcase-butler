package db

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	Module  string
	Process string
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
)

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer bytes.Buffer

	var levelColor string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = Green
	case logrus.WarnLevel:
		levelColor = Yellow
	case logrus.ErrorLevel:
		levelColor = Red
	case logrus.DebugLevel:
		levelColor = Blue
	default:
		levelColor = Reset
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	buffer.WriteString(fmt.Sprintf("[%s] ", timestamp))
	buffer.WriteString(fmt.Sprintf("%s[%s] ", levelColor, entry.Level.String()))
	if f.Module != "" && f.Process != "" {
		buffer.WriteString(fmt.Sprintf("[%s]:[%s] ", f.Module, f.Process))
	} else if f.Module != "" {
		buffer.WriteString(fmt.Sprintf("[%s] ", f.Module))
	}
	buffer.WriteString(entry.Message)
	buffer.WriteString(Reset)
	buffer.WriteString("\n")

	return buffer.Bytes(), nil
}

func InitLoggerStdout() {
	logrus.New()
	if CURRENT_DEBUG_STATUS {
		logrus.SetLevel(logrus.TraceLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}
