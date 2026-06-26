package notifier

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendAlert(webhookURL, message string) {
	if webhookURL == "" || webhookURL == "mock" {
		log.Printf("[MOCK DISCORD ALERT] %s\n", message)
		return
	}

	payload := map[string]string{
		"content": message,
	}
	data, _ := json.Marshal(payload)
	
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("⚠️ Failed to send discord alert: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		log.Printf("⚠️ Discord returned status code: %d. Check your webhook URL.\n", resp.StatusCode)
	} else {
		log.Println("✅ Discord alert sent successfully.")
	}
}
