package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 Starting Rubika Bot with Keyboard...")
	
	bot := NewRobot("BOT_TOKEN",
		WithTimeout(30*time.Second),
		WithPlatform("android"),
	)

	bot.OnMessage(func(r *Robot, m *Message) {
		fmt.Printf("📩 Received message from %s: %s\n", m.SenderID, m.Text)
		
		switch strings.TrimSpace(m.Text) {
		case "/start":
			fmt.Println("✅ Processing /start command")
			sendMainKeyboard(r, m.ChatID)
			
		case "📊 اطلاعات ربات":
			r.SendMessage(m.ChatID, "🤖 *اطلاعات ربات:*\n\n• نام: ربات تست\n• نسخه: 1.0.0\n• حالت: Polling\n• زبان: Go", nil)
			
		case "⭐ امتیازدهی":
			sendRatingKeyboard(r, m.ChatID)
			
		case "📞 تماس با پشتیبانی":
			r.SendMessage(m.ChatID, "📞 *پشتیبانی:*\n\n• ایدی: @Daniyel_Support\n• ایمیل: support@daniyel.ir\n• ساعت کاری: 9-17", nil)
			
		case "📍 موقعیت مکانی":
			r.SendLocation(m.ChatID, 35.6892, 51.3890, map[string]interface{}{
				"text": "📍 دفتر مرکزی - تهران",
			})
			
		case "🎵 ارسال موزیک":
			r.SendMessage(m.ChatID, "🎵 لطفاً یک فایل موزیک ارسال کنید...", nil)
			
		case "🖼 ارسال عکس":
			r.SendMessage(m.ChatID, "🖼 لطفاً یک عکس ارسال کنید...", nil)
			
		case "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐":
			handleRating(m.Text, r, m.ChatID)
			
		case "🔙 برگشت به منوی اصلی":
			sendMainKeyboard(r, m.ChatID)
			
		default:
			if strings.HasPrefix(m.Text, "/") {
				r.SendMessage(m.ChatID, "⚠️ دستور نامعتبر! از /start استفاده کنید.", nil)
			} else {
				r.SendMessage(m.ChatID, fmt.Sprintf("📨 شما گفتید: \"%s\"\n\n💡 از کیبورد پایین استفاده کنید.", m.Text), nil)
			}
		}
	})

	fmt.Println("⏳ Bot with keyboard is running...")
	fmt.Println("📩 Send /start to your bot")
	bot.Run()
}

func sendMainKeyboard(r *Robot, chatID string) {
	keyboard := map[string]interface{}{
		"rows": []map[string]interface{}{
			{
				"buttons": []map[string]interface{}{
					{"id": "info_btn", "type": "Simple", "button_text": "📊 اطلاعات ربات"},
					{"id": "rating_btn", "type": "Simple", "button_text": "⭐ امتیازدهی"},
				},
			},
			{
				"buttons": []map[string]interface{}{
					{"id": "contact_btn", "type": "Simple", "button_text": "📞 تماس با پشتیبانی"},
					{"id": "location_btn", "type": "Simple", "button_text": "📍 موقعیت مکانی"},
				},
			},
			{
				"buttons": []map[string]interface{}{
					{"id": "music_btn", "type": "Simple", "button_text": "🎵 ارسال موزیک"},
					{"id": "photo_btn", "type": "Simple", "button_text": "🖼 ارسال عکس"},
				},
			},
		},
		"resize_keyboard": true,
		"on_time_keyboard": false,
	}
	
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    "🎛 *منوی اصلی ربات*\n\nلطفاً یکی از گزینه‌های زیر را انتخاب کنید:",
		"chat_keypad": keyboard,
		"chat_keypad_type": "New",
	}
	
	_, err := r.post("sendMessage", payload)
	if err != nil {
		fmt.Printf("❌ Error sending keyboard: %v\n", err)
	} else {
		fmt.Println("✅ Main keyboard sent successfully")
	}
}

func sendRatingKeyboard(r *Robot, chatID string) {
	keyboard := map[string]interface{}{
		"rows": []map[string]interface{}{
			{
				"buttons": []map[string]interface{}{
					{"id": "star1", "type": "Simple", "button_text": "⭐"},
					{"id": "star2", "type": "Simple", "button_text": "⭐⭐"},
					{"id": "star3", "type": "Simple", "button_text": "⭐⭐⭐"},
				},
			},
			{
				"buttons": []map[string]interface{}{
					{"id": "star4", "type": "Simple", "button_text": "⭐⭐⭐⭐"},
					{"id": "star5", "type": "Simple", "button_text": "⭐⭐⭐⭐⭐"},
				},
			},
			{
				"buttons": []map[string]interface{}{
					{"id": "back_btn", "type": "Simple", "button_text": "🔙 برگشت به منوی اصلی"},
				},
			},
		},
		"resize_keyboard": true,
	}
	
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    "⭐ لطفاً به ربات امتیاز دهید:",
		"chat_keypad": keyboard,
		"chat_keypad_type": "New",
	}
	
	_, err := r.post("sendMessage", payload)
	if err != nil {
		fmt.Printf("❌ Error sending rating keyboard: %v\n", err)
	} else {
		fmt.Println("✅ Rating keyboard sent successfully")
	}
}

func handleRating(rating string, r *Robot, chatID string) {
	var response string
	switch rating {
	case "⭐":
		response = "😢 امتیاز 1 - متاسفم که ربات رو دوست نداشتی! چه چیزی رو می‌تونی بهتر کنیم؟"
	case "⭐⭐":
		response = "😐 امتیاز 2 - ممنون از بازخوردت! چه پیشنهادی داری؟"
	case "⭐⭐⭐":
		response = "😊 امتیاز 3 - ممنون! سعی می‌کنیم بهتر بشیم."
	case "⭐⭐⭐⭐":
		response = "😄 امتیاز 4 - عالی! خوشحالیم که ربات رو دوست داری."
	case "⭐⭐⭐⭐⭐":
		response = "🎉 امتیاز 5 - فوق‌العاده! ممنون از انرژی مثبتت 💫"
	default:
		response = "⚠️ امتیاز نامعتبر"
	}
	
	r.SendMessage(chatID, response, nil)
	time.Sleep(2 * time.Second)
	sendMainKeyboard(r, chatID)
}