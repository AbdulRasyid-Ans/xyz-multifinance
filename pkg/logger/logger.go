package logger

import (
	"log"
	"os"
)

func Info(message string) {
	infoLogger := log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime)
	infoLogger.Println(message)
}

func Error(message string) {
	errorLogger := log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime)
	errorLogger.Println(message)
}

func Warning(message string) {
	warningLogger := log.New(os.Stdout, "[WARNING]: ", log.Ldate|log.Ltime)
	warningLogger.Println(message)
}
