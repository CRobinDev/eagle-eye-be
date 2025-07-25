package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

func NewLogger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetReportCaller(true)

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "Jan 02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				repopath, _ := os.Getwd()
				repopath = filepath.ToSlash(repopath) + "/"
				filename := strings.Replace(f.File, repopath, " ", 1)
				return "", fmt.Sprintf("%s:%d", filename, f.Line)
			},
			ForceColors:      true,
			QuoteEmptyFields: true,
		})

		repoPath, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("failed to get current working directory: %v", err))
		}

		logFileName := filepath.Join(repoPath, "/storage/logs", fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open log file: %v", err))
		}

		multiWriter := io.MultiWriter(os.Stdout, file)
		logger.SetOutput(multiWriter)

		if env.GetEnv().AppEnv == "production" {
			logger.SetLevel(logrus.InfoLevel)
		} else {
			logger.SetLevel(logrus.DebugLevel)
		}
	})

	return logger
}

func WithTraceID(traceID uuid.UUID, err error) logrus.Fields {
	if err == nil {
		return logrus.Fields{
			"trace_id": traceID,
		}
	} else {
		return logrus.Fields{
			"error":    err,
			"trace_id": traceID,
		}
	}
}
