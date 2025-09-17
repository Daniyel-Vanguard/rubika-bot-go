package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 Starting Rubika Bot in Polling Mode...")
	
	bot := NewRobot("BOT_TOKEN",
		WithTimeout(30*time.Second),
		WithPlatform("android"),
	)

	bot.OnMessage(func(r *Robot, m *Message) {
		fmt.Printf("📩 Received message from %s: %s\n", m.SenderID, m.Text)
		
		if strings.HasPrefix(m.Text, "/start") {
			fmt.Println("✅ Processing /start command")
			
			btn1 := CreateInlineButton("📊 اطلاعات ربات", "bot_info", "Simple")
			btn2 := CreateInlineButton("ℹ️ راهنما", "help", "Simple")
			
			row := CreateButtonRow(btn1, btn2)
			
			keypad := CreateInlineKeypad([]map[string]interface{}{row})
			
			payload := map[string]interface{}{
				"chat_id": m.ChatID,
				"text":    "🤖 به ربات خوش آمدید!\n\nبرای شروع از دکمه‌های زیر استفاده کنید:",
				"inline_keypad": keypad,
			}
			
			result, err := r.post("sendMessage", payload)
			
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Println("✅ Message sent, checking response...")
				if result != nil {
					jsonData, _ := json.MarshalIndent(result, "", "  ")
					fmt.Printf("📋 API Response: %s\n", string(jsonData))
					
					if status, ok := result["status"].(string); ok {
						fmt.Printf("📊 API Status: %s\n", status)
						if status != "OK" && status != "ok" {
							if data, ok := result["data"].(map[string]interface{}); ok {
								if errorMsg, ok := data["error_message"].(string); ok {
									fmt.Printf("❌ Error Message: %s\n", errorMsg)
								}
							}
						}
					}
				}
			}
		}
	})

	bot.OnCallback("bot_info", func(r *Robot, m *Message) {
		fmt.Println("✅ Button bot_info clicked")
		r.SendMessage(m.ChatID, "📊 اطلاعات ربات: این یک ربات تست است", nil)
	})

	bot.OnCallback("help", func(r *Robot, m *Message) {
		fmt.Println("✅ Button help clicked")
		r.SendMessage(m.ChatID, "ℹ️ راهنما: از /start استفاده کنید", nil)
	})

	fmt.Println("⏳ Bot is running...")
	bot.Run()
}
