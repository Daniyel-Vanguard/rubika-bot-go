package main

import (
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
		
		switch {
		case strings.HasPrefix(m.Text, "/start"):
			btn1 := CreateInlineButton("📊 اطلاعات ربات", "bot_info", "Primary")
			btn2 := CreateInlineButton("ℹ️ راهنما", "help", "Secondary")
			row := []map[string]interface{}{btn1, btn2}
			keypad := CreateInlineKeypad(row)
			
			r.SendMessage(m.ChatID, "🤖 به ربات خوش آمدید!\n\nبرای شروع از دکمه‌های زیر استفاده کنید:", map[string]interface{}{
				"inline_keypad": keypad,
			})

		case strings.HasPrefix(m.Text, "/help"):
			helpText := `📖 راهنمای ربات:
/start - شروع کار با ربات
/help - نمایش راهنما
/ping - تست ارتباط`
			
			r.SendMessage(m.ChatID, helpText)

		case strings.HasPrefix(m.Text, "/ping"):
			r.SendMessage(m.ChatID, "🏓 Pong! ربات فعال است!")

		default:
			r.SendMessage(m.ChatID, fmt.Sprintf("👤 شما گفتید: %s\n\n💡 از /help برای راهنما استفاده کنید.", m.Text))
		}
	})

	bot.OnCallback("bot_info", func(r *Robot, m *Message) {
		info, err := r.GetMe()
		if err == nil {
			botName := "ربات"
			if data, ok := info["data"].(map[string]interface{}); ok {
				if botData, ok := data["bot"].(map[string]interface{}); ok {
					if name, ok := botData["name"].(string); ok {
						botName = name
					}
				}
			}
			r.SendMessage(m.ChatID, fmt.Sprintf("🤖 نام ربات: %s\n✅ وضعیت: فعال\n🎯 حالت: Polling", botName))
		}
	})

	bot.OnCallback("help", func(r *Robot, m *Message) {
		helpText := `🎯 دستورات قابل استفاده:
• /start - شروع کار
• /help - راهنمایی
• /ping - تست ارتباط

🤖 این ربات در حالت Polling کار می‌کند.`
		
		r.SendMessage(m.ChatID, helpText)
	})

	fmt.Println("⏳ Bot is running in polling mode...")
	bot.Run()
}
