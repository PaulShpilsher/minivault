package infrastructure

import (
	"encoding/json"
	"os"
	"time"
	"log"
	"minivault/domain"
)

const LOG_FILE = "logs/log.jsonl"

func LogInteraction(input domain.GenerateRequest, output domain.GenerateResponse) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"input":     input,
		"output":    output,
	}
	logData, _ := json.Marshal(logEntry)
	os.MkdirAll("logs", 0755)
	f, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open log file: %v", err)
		return
	}
	defer f.Close()
	f.Write(append(logData, '\n'))
}
