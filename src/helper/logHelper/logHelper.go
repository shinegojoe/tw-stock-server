package logHelper

import (
	"log"
	"os"
)

func LogToFile(message string) {
	// cfg := config.GetConfig()
	f, err := os.OpenFile("log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "error ", log.LstdFlags)
	logger.Println(message)
}
