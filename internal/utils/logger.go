package utils

import (
	"log"
	"os"
)

// LogError logs an error to console and file
func LogError(err error) {
	log.Println("ERROR:", err)
	f, _ := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	log.SetOutput(f)
	log.Println("ERROR:", err)
}