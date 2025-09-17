# 🤖 کتابخانه روبیکا برای Go - Rubika Bot Go Library
  <img align="center" src="https://rubika.ir/static/images/logo.svg"/>
  <br/>
  <img align="center" src="https://go.dev/images/gophers/ladder.svg"/>

یک کتابخانه ساده و قدرتمند برای ساخت ربات‌های روبیکا با زبان Go

# 📦 نصب

```bash
git clone https://github.com/Daniyel-Vanguard/rubika-bot-go.git
```

# 🚀 شروع سریع

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // ایجاد ربات با توکن
    bot := rubika.NewRobot("YOUR_BOT_TOKEN_HERE",
        rubika.WithTimeout(30*time.Second),
        rubika.WithPlatform("android"),
    )

    // ثبت هندلر برای پیام‌ها
    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        if m.Text == "/start" {
            // ایجاد کیبورد ساده
            keyboard := map[string]interface{}{
                "rows": []map[string]interface{}{
                    {
                        "buttons": []map[string]interface{}{
                            {"id": "btn1", "type": "Simple", "button_text": "📊 اطلاعات"},
                            {"id": "btn2", "type": "Simple", "button_text": "⭐ امتیاز"},
                        },
                    },
                },
                "resize_keyboard": true,
            }
            
            r.SendMessage(m.ChatID, "سلام! به ربات خوش آمدید 👋", map[string]interface{}{
                "chat_keypad": keyboard,
                "chat_keypad_type": "New",
            })
        }
    })

    // شروع ربات
    fmt.Println("🤖 ربات در حال اجرا...")
    bot.Run()
}
```

# 📖 مستندات کامل

📋 ساختار اصلی

```go
// ایجاد ربات جدید
bot := rubika.NewRobot(token string, options ...func(*Robot))

// آپشن‌های قابل تنظیم
```

💬 مدیریت پیام‌ها

```go
// ثبت هندلر برای پیام‌های متنی
bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
    fmt.Printf("پیام از %s: %s\n", m.SenderID, m.Text)
    
    // پاسخ به پیام
    r.SendMessage(m.ChatID, "پیام شما دریافت شد!", nil)
})

// ثبت هندلر برای callback دکمه‌ها
bot.OnCallback("button_id", func(r *rubika.Robot, m *rubika.Message) {
    r.SendMessage(m.ChatID, "دکمه کلیک شد!", nil)
})
```

⌨️ ایجاد کیبورد

```go
// ایجاد کیبورد معمولی (چت کیپد)
keyboard := map[string]interface{}{
    "rows": []map[string]interface{}{
        {
            "buttons": []map[string]interface{}{
                {
                    "id": "btn1",
                    "type": "Simple", 
                    "button_text": "دکمه ۱"
                },
                {
                    "id": "btn2",
                    "type": "Simple",
                    "button_text": "دکمه ۲"
                },
            },
        },
    },
    "resize_keyboard": true,
}

// ارسال پیام با کیبورد
r.SendMessage(chatID, "پیام با کیبورد", map[string]interface{}{
    "chat_keypad": keyboard,
    "chat_keypad_type": "New",
})
```

📤 ارسال انواع محتوا

```go
// ارسال متن ساده
r.SendMessage(chatID, "سلام دنیا!", nil)

// ارسال موقعیت مکانی
r.SendLocation(chatID, 35.6892, 51.3890, map[string]interface{}{
    "text": "این موقعیت من است",
})

// ارسال مخاطب
r.SendContact(chatID, "جان", "دو", "+989123456789", nil)

// ارسال نظرسنجی
r.SendPoll(chatID, "نظر شما چیست؟", []string{"گزینه ۱", "گزینه ۲", "گزینه ۳"})
```

🖼️ ارسال فایل و مدیا

```go
// ارسال عکس
r.SendImage(chatID, "path/to/image.jpg", map[string]interface{}{
    "text": "این یک عکس است",
})

// ارسال فایل
r.SendDocument(chatID, "path/to/file.pdf", map[string]interface{}{
    "text": "این یک فایل است",
})

// ارسال موزیک
r.SendMusic(chatID, "path/to/music.mp3", nil)

// ارسال صوت
r.SendVoice(chatID, "path/to/voice.ogg", nil)

// ارسال GIF
r.SendGif(chatID, "path/to/animation.gif", nil)
```

# 🎯 مثال‌های کاربردی

مثال ۱: ربات پرسش و پاسخ

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    bot := rubika.NewRobot("YOUR_TOKEN",
        rubika.WithTimeout(30*time.Second),
    )

    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        switch m.Text {
        case "/start":
            sendWelcomeMenu(r, m.ChatID)
        case "📊 اطلاعات":
            r.SendMessage(m.ChatID, "🤖 این یک ربات نمونه است", nil)
        case "⭐ امتیاز":
            sendRatingMenu(r, m.ChatID)
        case "📞 پشتیبانی":
            r.SendMessage(m.ChatID, "📞 برای پشتیبانی با @Support联系 کنید", nil)
        default:
            r.SendMessage(m.ChatID, "⚠️ دستور نامعتبر!", nil)
        }
    })

    fmt.Println("🤖 ربات در حال اجرا...")
    bot.Run()
}

func sendWelcomeMenu(r *rubika.Robot, chatID string) {
    keyboard := map[string]interface{}{
        "rows": []map[string]interface{}{
            {
                "buttons": []map[string]interface{}{
                    {"id": "info", "type": "Simple", "button_text": "📊 اطلاعات"},
                    {"id": "rate", "type": "Simple", "button_text": "⭐ امتیاز"},
                },
            },
            {
                "buttons": []map[string]interface{}{
                    {"id": "support", "type": "Simple", "button_text": "📞 پشتیبانی"},
                },
            },
        },
        "resize_keyboard": true,
    }
    
    r.SendMessage(chatID, "🎉 به ربات خوش آمدید!", map[string]interface{}{
        "chat_keypad": keyboard,
        "chat_keypad_type": "New",
    })
}
```

مثال ۲: ربات نظرسنجی

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    bot := rubika.NewRobot("YOUR_TOKEN")

    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        if m.Text == "/start" {
            // ایجاد نظرسنجی
            question := "به چه زبانی برنامه‌نویسی می‌کنید؟"
            options := []string{"Go", "Python", "JavaScript", "Java", "C++"}
            
            r.SendPoll(m.ChatID, question, options)
        }
    })

    bot.Run()
}
```

# 🔧 مدیریت خطاها

```go
bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
    result, err := r.SendMessage(m.ChatID, "پیام تست", nil)
    if err != nil {
        fmt.Printf("❌ خطا در ارسال پیام: %v\n", err)
        return
    }
    
    // بررسی وضعیت پاسخ
    if status, ok := result["status"].(string); ok {
        if status != "OK" {
            fmt.Printf("⚠️ وضعیت API: %s\n", status)
        }
    }
})
```

# 🌐 وب‌هوک (اختیاری)

```go
// راه‌اندازی ربات با وب‌هوک
bot := rubika.NewRobot("YOUR_TOKEN",
    rubika.WithWebhook("https://yourdomain.com/webhook"),
)

// راه‌اندازی سرور وب‌هوک
err := bot.StartWebhookServer("8080")
if err != nil {
    log.Fatal("خطا در راه‌اندازی وب‌هوک: ", err)
}
```

# 📊 لاگ و دیباگ

```go
// فعال کردن لاگ پیشرفته
bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
    fmt.Printf("👤 کاربر: %s\n", m.SenderID)
    fmt.Printf("📩 پیام: %s\n", m.Text)
    fmt.Printf("💬 چت: %s\n", m.ChatID)
    fmt.Printf("📦 داده خام: %+v\n", m.RawData)
})
```

# 🚀 استقرار

اجرای本地

```bash
go run rubika_bot.go YOUR_BOT.go
```

# 🤝 مشارکت

1. فورک ریپو
2. ایجاد برنچ جدید
3. commit تغییرات
4. push به برنچ
5. ایجاد Pull Request

# 📜 لایسنس

این پروژه تحت لایسنس MIT منتشر شده است.

# 📞 پشتیبانی

· 📧 ایمیل: hadipishghadam13@gmail.com
· 🐛 issues: GitHub Issues

---

⭐ اگر این پروژه رو دوست داشتید، ستاره بدید!
