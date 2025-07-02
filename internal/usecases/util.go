package usecases

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func sendWebhook(userId, ip string) {
	payload := map[string]string{
		"user_id":    userId,
		"ip_address": ip,
		"timestamp":  time.Now().String(),
	}
	body, _ := json.Marshal(payload)
	http.Post("http://example.com/webhook", "application/json", strings.NewReader(string(body)))
}
