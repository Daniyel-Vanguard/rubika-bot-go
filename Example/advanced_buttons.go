package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 Starting Rubika Bot with Advanced Buttons...")
	
	bot := NewRobot("BOT_TOKEN",
		WithTimeout(30*time.Second),
		WithPlatform("android"),
	)

	bot.OnMessage(func(r *Robot, m *Message) {
		fmt.Printf("📩 Received message from %s: %s\n", m.SenderID, m.Text)
		
		if strings.HasPrefix(m.Text, "/start") {
			fmt.Println("✅ Processing /start command")
			
			btn1 := CreateInlineButton("📊 اطلاعات", "btn_info", "Simple")
			btn2 := CreateInlineButton("⭐ امتیازدهی", "btn_rating", "Simple")
			btn3 := CreateInlineButton("📞 تماس", "btn_contact", "Simple")
			btn4 := CreateInlineButton("📍 موقعیت", "btn_location", "Simple")
			btn5 := CreateInlineButton("🎵 موزیک", "btn_music", "Simple")
			btn6 := CreateInlineButton("🖼 عکس", "btn_photo", "Simple")
			
			row1 := CreateButtonRow(btn1, btn2)
			row2 := CreateButtonRow(btn3, btn4)
			row3 := CreateButtonRow(btn5, btn6)
			
			keypad := CreateInlineKeypad([]map[string]interface{}{row1, row2, row3})
			
			payload := map[string]interface{}{
				"chat_id": m.ChatID,
				"text":    "🎛 *منوی اصلی ربات*\n\nلطفاً یکی از گزینه‌های زیر را انتخاب کنید:",
				"inline_keypad": keypad,
			}
			
			_, err := r.post("sendMessage", payload)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Println("✅ Advanced buttons sent successfully")
			}

		} else if strings.HasPrefix(m.Text, "/test") {
			advancedKeypad := map[string]interface{}{
				"rows": []map[string]interface{}{
					{
						"buttons": []map[string]interface{}{
							{
								"id":          "camera_btn",
								"type":        "CameraImage",
								"button_text": "📷 دوربین",
							},
							{
								"id":          "gallery_btn",
								"type":        "GalleryImage", 
								"button_text": "🖼 گالری",
							},
						},
					},
					{
						"buttons": []map[string]interface{}{
							{
								"id":          "location_btn",
								"type":        "MyLocation",
								"button_text": "📍 موقعیت من",
							},
							{
								"id":          "phone_btn",
								"type":        "MyPhoneNumber",
								"button_text": "📞 شماره من",
							},
						},
					},
					{
						"buttons": []map[string]interface{}{
							{
								"id":          "audio_btn",
								"type":        "Audio",
								"button_text": "🎵 ارسال صوت",
							},
							{
								"id":          "file_btn",
								"type":        "File",
								"button_text": "📁 ارسال فایل",
							},
						},
					},
				},
			}
			
			payload := map[string]interface{}{
				"chat_id": m.ChatID,
				"text":    "📋 *منوی پیشرفته*\n\nاین دکمه‌های خاصیت‌های مختلفی دارند:",
				"inline_keypad": advancedKeypad,
			}
			
			_, err := r.post("sendMessage", payload)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				fmt.Println("✅ Special buttons sent successfully")
			}
		}
	})

	bot.OnCallback("btn_info", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "🤖 *اطلاعات ربات:*\n\n• نام: ربات تست\n• نسخه: 1.0.0\n• حالت: Polling", nil)
	})

	bot.OnCallback("btn_rating", func(r *Robot, m *Message) {
		star1 := CreateInlineButton("⭐", "star_1", "Simple")
		star2 := CreateInlineButton("⭐⭐", "star_2", "Simple")
		star3 := CreateInlineButton("⭐⭐⭐", "star_3", "Simple")
		star4 := CreateInlineButton("⭐⭐⭐⭐", "star_4", "Simple")
		star5 := CreateInlineButton("⭐⭐⭐⭐⭐", "star_5", "Simple")
		
		ratingRow := CreateButtonRow(star1, star2, star3, star4, star5)
		keypad := CreateInlineKeypad([]map[string]interface{}{ratingRow})
		
		r.SendMessage(m.ChatID, "⭐ لطفاً به ربات امتیاز دهید:", map[string]interface{}{
			"inline_keypad": keypad,
		})
	})

	bot.OnCallback("btn_contact", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "📞 برای تماس با پشتیبانی:\n@Daniyel_Support", nil)
	})

	bot.OnCallback("btn_location", func(r *Robot, m *Message) {
		r.SendLocation(m.ChatID, 35.6892, 51.3890, map[string]interface{}{
			"text": "📍 موقعیت دفمر مرکزی",
		})
	})

	bot.OnCallback("btn_music", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "🎵 این قابلیت به زودی اضافه خواهد شد...", nil)
	})

	bot.OnCallback("btn_photo", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "🖼 این قابلیت به زودی اضافه خواهد شد...", nil)
	})

	bot.OnCallback("camera_btn", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "📷 دسترسی به دوربین باز شد...", nil)
	})

	bot.OnCallback("gallery_btn", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "🖼 دسترسی به گالری باز شد...", nil)
	})

	bot.OnCallback("location_btn", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "📍 موقعیت شما دریافت شد!", nil)
	})

	bot.OnCallback("phone_btn", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "📞 شماره تلفن شما دریافت شد!", nil)
	})

	bot.OnCallback("star_1", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "😢 متاسفم که ربات رو دوست نداشتی! چه چیزی رو می‌تونی بهتر کنیم؟", nil)
	})

	bot.OnCallback("star_5", func(r *Robot, m *Message) {
		r.SendMessage(m.ChatID, "🎉 ممنون از امتیاز عالیت! خوشحالیم که ربات رو دوست داری!", nil)
	})

	fmt.Println("⏳ Bot with advanced buttons is running...")
	fmt.Println("📩 Send /start or /test to your bot")
	bot.Run()
}