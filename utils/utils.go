package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Any = interface{}

func Use(vals ...interface{}) {}

func Perform[A Any](cont func() A) chan A {
	ch := make(chan A)
	go func() {
		ch <- cont()
	}()
	return ch
}

func Try(cont func() interface{}) chan interface{} {
	ch := make(chan interface{})
	go func() {
		ch <- cont()
	}()
	return ch
}

func InitLogger() error {
	loggingFile := os.Getenv("LOG_FILE")
	file, err := os.OpenFile(loggingFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
