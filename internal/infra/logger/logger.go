package logger

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	instance *log.Logger
	once     sync.Once
	logFile  *os.File
)

func Init(logFilePath string) error {
	var err error
	once.Do(func() {
		initLogger(logFilePath, &err)
	})
	return err
}

func initLogger(logFilePath string, err *error) {
	var output io.Writer = os.Stdout

	if logFilePath != "" {
		logFile, *err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if *err == nil {
			output = logFile
		}
	}

	instance = log.New(output, "", log.LstdFlags|log.Lshortfile)

	log.SetOutput(output)
}

func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func Debug(v ...interface{}) {
	ensureInitialized()
	instance.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	ensureInitialized()
	instance.Printf(format, v...)
}

func Info(v ...interface{}) {
	ensureInitialized()
	instance.Println(v...)
}

func Infof(format string, v ...interface{}) {
	ensureInitialized()
	instance.Printf(format, v...)
}

func Warn(v ...interface{}) {
	ensureInitialized()
	instance.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	ensureInitialized()
	instance.Printf(format, v...)
}

func Error(v ...interface{}) {
	ensureInitialized()
	instance.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	ensureInitialized()
	instance.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	ensureInitialized()
	instance.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	ensureInitialized()
	instance.Fatalf(format, v...)
}

func ensureInitialized() {
	if instance == nil {
		var err error
		once.Do(func() {
			initLogger("", &err)
		})

		if instance == nil {
			instance = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		}
	}
}
