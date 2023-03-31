package file

import (
	"log"
	"os"
)

type FileLogger struct{}

func (self *FileLogger) Info(msg string) {
	log.Println("[Info] " + msg)
}
func (self *FileLogger) Error(err error) {
	log.Println("[Error] " + err.Error())
}
func (self *FileLogger) Debug(s string) {
	log.Println("[Debug] " + s)
}
func (self *FileLogger) Fatal(err error) {
	log.Fatal("[Fatal] " + err.Error())
}

func initLogger() error {
	loggingFile := os.Getenv("LOG_FILE")
	file, err := os.OpenFile(loggingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}

func Instance() (FileLogger, error) {
	logger := FileLogger{}
	err := initLogger()
	if err != nil {
		return logger, err
	}
	return logger, nil
}
