package services

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func NewLogger(infoWriter, warningWriter, errorWriter io.Writer) *Logger {
	return &Logger{
		InfoLogger:    log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLogger: log.New(warningWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger:   log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func InitLogger() *Logger {
	infoLogFile, err := os.OpenFile("/logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening info log file:", err)
		return nil
	}
	warningLogFile, err := os.OpenFile("/logs/warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening warning log file:", err)
		return nil
	}
	errorLogFile, err := os.OpenFile("/logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening error log file:", err)
		return nil
	}

	return NewLogger(io.MultiWriter(os.Stdout, infoLogFile), io.MultiWriter(os.Stdout, warningLogFile), io.MultiWriter(os.Stdout, errorLogFile))
}
