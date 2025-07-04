package usecases

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func sendWebhook(userId, ip string) {
	message := map[string]string{
		"user_id":    userId,
		"ip_address": ip,
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	}
	body, _ := json.Marshal(message)
	http.Post("http://localhost:8080/auth/webhook", "application/json", strings.NewReader(string(body)))
}
