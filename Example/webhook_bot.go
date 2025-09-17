package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 Starting Rubika Bot in Webhook Mode...")
	
	bot := NewRobot("BOT_TOKEN",
		WithWebhook("https://yourdomain.com:8080/webhook"),
		WithTimeout(30*time.Second),
		WithPlatform("android"),
	)

	bot.OnMessage(func(r *Robot, m *Message) {
		fmt.Printf("🌐 Webhook received from %s: %s\n", m.SenderID, m.Text)
		
		command := strings.ToLower(strings.TrimSpace(m.Text))
		switch command {
		case "/start", "start":
			welcomeMsg := `🎉 به ربات وب‌هوک خوش آمدید!
			
این ربات از طریق وب‌هوک پیام‌ها را دریافت می‌کند.

دستورات موجود:
/help - نمایش راهنما
/ping - تست ارتباط
/info - اطلاعات ربات`
			
			r.SendMessage(m.ChatID, welcomeMsg)

		case "/help", "help":
			helpMsg := `📖 راهنمای ربات وب‌هوک:

• این ربات با تکنولوژی Webhook کار می‌کند
• پیام‌ها مستقیماً از سرور روبیکا دریافت می‌شوند
• سرعت پاسخگویی بالاتر از حالت Polling است
• مناسب برای ربات‌های پرترافیک`

			r.SendMessage(m.ChatID, helpMsg)

		case "/ping", "ping":
			r.SendMessage(m.ChatID, "🏓 Pong! Connection is working perfectly!")

		case "/info", "info":
			botInfo, err := r.GetMe()
			if err == nil {
				infoText := "🤖 اطلاعات ربات:\n"
				if data, ok := botInfo["data"].(map[string]interface{}); ok {
					if botData, ok := data["bot"].(map[string]interface{}); ok {
						if name, ok := botData["name"].(string); ok {
							infoText += fmt.Sprintf("نام: %s\n", name)
						}
						if username, ok := botData["username"].(string); ok {
							infoText += fmt.Sprintf("آیدی: @%s\n", username)
						}
					}
				}
				infoText += "حالت: Webhook 🌐"
				r.SendMessage(m.ChatID, infoText)
			}

		default:
			echoMsg := fmt.Sprintf("📨 پیام شما: %s\n\n💡 از /help برای راهنما استفاده کنید.", m.Text)
			r.SendMessage(m.ChatID, echoMsg)
		}
	})

	fmt.Println("🌐 Setting up webhook server...")
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "healthy",
			"bot":       "running",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		info, _ := bot.GetMe()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "online",
			"mode":        "webhook",
			"bot_info":    info,
			"webhook_url": bot.WebhookURL,
		})
	})

	fmt.Println("✅ Webhook server starting on port 8080...")
	fmt.Println("📊 Health check: http://localhost:8080/health")
	fmt.Println("📈 Status: http://localhost:8080/status")
	
	if err := bot.StartWebhookServer("8080"); err != nil {
		log.Fatalf("❌ Failed to start webhook server: %v", err)
	}
}
