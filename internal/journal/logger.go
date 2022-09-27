package journal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

func New(part string, writeFile bool) {
	once.Do(func() {
		logger = logrus.New()

		if writeFile {
			filePath := filepath.Join("logs", part, " ", time.Now().Format("2006-01-02 15:04:05"))
			file, err := os.OpenFile(filePath, os.O_CREATE, os.ModeAppend)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Logger"))
			}
			writer := io.Writer(file)
			logger.SetOutput(writer)
		}

		logger.SetLevel(logrus.InfoLevel)

		logger.SetFormatter(
			&logrus.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			},
		)
	})
}

func LogInfo(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo()
	entry.Info(args...)
}

func LogWarn(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo()
	entry.Warn(args...)
}

func LogError(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo()
	entry.Error(args...)
}

func LogErrorf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo()
	entry.Errorf(format, args...)
}

func LogFatal(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	entry.Data["file"] = fileInfo()
	entry.Fatal(args...)
}

func fileInfo() string {
	_, file, line, ok := runtime.Caller(2) // skip 2 calls from stack for fileInfo() and Log... functions to get real caller
	if !ok {
		file = "<???>"
		line = 1
	}

	return fmt.Sprintf("%s:%d", file, line)
}
