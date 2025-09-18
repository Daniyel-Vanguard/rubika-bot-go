# 📚 مستندات کامل کتابخانه Rubika Bot برای Go

# 🎯 معرفی

یک کتابخانه جامع و قدرتمند برای ساخت ربات‌های پیام‌رسان روبیکا با زبان Go. این کتابخانه تمامی امکانات API روبیکا را پوشش می‌دهد و ساخت ربات‌های حرفه‌ای را بسیار ساده می‌کند.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    bot := rubika.NewRobot("YOUR_BOT_TOKEN")
    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        r.SendMessage(m.ChatID, "سلام! 👋", nil)
    })
    bot.Run()
}
```

# 📦 نصب و راه‌اندازی

پیش‌نیازها

· Go 1.16 یا بالاتر
· توکن ربات 

نصب

```bash
git clone https://github.com/Daniyel-Vanguard/rubika-bot-go.git
```

# 🔧 پیکربندی اولیه

ایجاد نمونه ربات

```go
// ساده‌ترین روش
bot := rubika.NewRobot("YOUR_BOT_TOKEN")

// با تنظیمات پیشرفت
)
```

# 💬 مدیریت پیام‌ها

هندلر پیام متنی

```go
bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
    fmt.Printf("📨 پیام جدید از %s: %s\n", m.SenderID, m.Text)
    
    // پاسخ به پیام
    r.SendMessage(m.ChatID, "پیام شما دریافت شد!", nil)
})
```

ساختار شیء Message

```go
type Message struct {
    Bot       *Robot
    ChatID    string                 // آیدی چت
    MessageID string                 // آیدی پیام
    SenderID  string                 // آیدی فرستنده
    Text      string                 // متن پیام
    RawData   map[string]interface{} // داده خام
}
```

پردازش دستورات

```go
bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
    switch m.Text {
    case "/start":
        handleStartCommand(r, m)
    case "/help":
        handleHelpCommand(r, m)
    case "/settings":
        handleSettingsCommand(r, m)
    default:
        handleOtherMessages(r, m)
    }
})

func handleStartCommand(r *rubika.Robot, m *rubika.Message) {
    welcomeText := `🎉 به ربات خوش آمدید!

دستورات موجود:
/start - شروع کار
/help - راهنما
/settings - تنظیمات`

    r.SendMessage(m.ChatID, welcomeText, nil)
}
```

# ⌨️ مدیریت کیبورد و دکمه‌ها

ایجاد کیبورد ساده

```go
func createSimpleKeyboard() map[string]interface{} {
    return map[string]interface{}{
        "rows": []map[string]interface{}{
            {
                "buttons": []map[string]interface{}{
                    {
                        "id": "btn_info",
                        "type": "Simple",
                        "button_text": "📊 اطلاعات",
                    },
                    {
                        "id": "btn_help", 
                        "type": "Simple",
                        "button_text": "ℹ️ راهنما",
                    },
                },
            },
        },
        "resize_keyboard": true,
    }
}
```

ارسال پیام با کیبورد

```go
func sendWelcomeMessage(r *rubika.Robot, chatID string) {
    keyboard := createSimpleKeyboard()
    
    payload := map[string]interface{}{
        "chat_id": chatID,
        "text": "لطفاً یک گزینه انتخاب کنید:",
        "chat_keypad": keyboard,
        "chat_keypad_type": "New",
    }
    
    r.SendMessage(chatID, "به ربات خوش آمدید!", payload)
}
```

انواع دکمه‌ها

```go
// دکمه ساده
btnSimple := map[string]interface{}{
    "id": "simple_btn",
    "type": "Simple",
    "button_text": "دکمه ساده",
}

// دکمه تماس
btnContact := map[string]interface{}{
    "id": "contact_btn", 
    "type": "MyPhoneNumber",
    "button_text": "📞 شماره من",
}

// دکمه موقعیت
btnLocation := map[string]interface{}{
    "id": "location_btn",
    "type": "MyLocation", 
    "button_text": "📍 موقعیت من",
}

// دکمه دوربین
btnCamera := map[string]interface{}{
    "id": "camera_btn",
    "type": "CameraImage",
    "button_text": "📷 دوربین",
}
```

# 📤 ارسال انواع محتوا

ارسال متن ساده

```go
// متن ساده
r.SendMessage(chatID, "سلام دنیا!", nil)

// متن با فرمت
r.SendMessage(chatID, "*متن ضخیم* و _متن کج_", nil)

// متن چندخطی
message := `خط اول
خط دوم
خط سوم`

r.SendMessage(chatID, message, nil)
```

ارسال موقعیت مکانی

```go
// ارسال موقعیت
r.SendLocation(chatID, 35.6892, 51.3890, map[string]interface{}{
    "text": "📍 این موقعیت من است",
    "disable_notification": false,
})
```

ارسال مخاطب

```go
// ارسال مخاطب
r.SendContact(chatID, "جان", "دو", "+989123456789", map[string]interface{}{
    "text": "👤 اطلاعات تماس",
})
```

ارسال نظرسنجی

```go
// ایجاد نظرسنجی
question := "نظر شما درباره این ربات چیست؟"
options := []string{"عالی", "خوب", "متوسط", "ضعیف"}

r.SendPoll(chatID, question, options)
```

# 🖼️ ارسال فایل و مدیا

ارسال عکس

```go
// ارسال عکس از مسیر本地
r.SendImage(chatID, "./images/welcome.jpg", map[string]interface{}{
    "text": "عکس خوش‌آمدگویی",
    "disable_notification": true,
})

// ارسال عکس از URL
r.SendImage(chatID, "https://example.com/image.jpg", nil)
```

ارسال فایل

```go
// ارسال سند
r.SendDocument(chatID, "./files/document.pdf", map[string]interface{}{
    "text": "📄 فایل PDF",
})

// ارسال فایل از URL
r.SendDocument(chatID, "https://example.com/file.zip", nil)
```

ارسال صوت و موزیک

```go
// ارسال صوت
r.SendVoice(chatID, "./audio/message.ogg", map[string]interface{}{
    "text": "🎵 پیام صوتی",
})

// ارسال موزیک
r.SendMusic(chatID, "./music/song.mp3", map[string]interface{}{
    "text": "🎶 آهنگ جدید",
})
```

ارسال GIF

```go
// ارسال GIF
r.SendGif(chatID, "./gifs/animation.gif", map[string]interface{}{
    "text": "🎬 انیمیشن GIF",
})
```

# 🔄 مدیریت وضعیت ربات

حالت Polling

```go
func main() {
    bot := rubika.NewRobot("YOUR_TOKEN")
    
    // تنظیم هندلرها
    bot.OnMessage(messageHandler)
    bot.OnCallback("btn_info", infoHandler)
    
    // شروع ربات
    fmt.Println("🤖 ربات در حال اجرا...")
    bot.Run()
}
```

حالت Webhook

```go
func main() {
    bot := rubika.NewRobot("YOUR_TOKEN",
        rubika.WithWebhook("https://yourdomain.com/webhook"),
    )
    
    // راه‌اندازی وب‌هوک
    err := bot.StartWebhookServer("8080")
    if err != nil {
        log.Fatal("خطا در راه‌اندازی وب‌هوک: ", err)
    }
    
    fmt.Println("🌐 وب‌هوک فعال روی پورت 8080")
}
```

🛠️ ابزارهای کمکی

ایجاد دکمه

```go
// ایجاد دکمه ساده
func createButton(text, buttonID, buttonType string) map[string]interface{} {
    return map[string]interface{}{
        "id": buttonID,
        "type": buttonType,
        "button_text": text,
    }
}

// ایجاد ردیف دکمه
func createButtonRow(buttons ...map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "buttons": buttons,
    }
}

// ایجاد کیبورد
func createKeyboard(rows ...map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "rows": rows,
        "resize_keyboard": true,
    }
}
```

مدیریت خطاها

```go
func safeSendMessage(r *rubika.Robot, chatID, text string) {
    result, err := r.SendMessage(chatID, text, nil)
    if err != nil {
        fmt.Printf("❌ خطا در ارسال پیام: %v\n", err)
        return
    }
    
    // بررسی وضعیت پاسخ
    if status, ok := result["status"].(string); ok {
        if status != "OK" {
            fmt.Printf("⚠️ وضعیت API: %s\n", status)
            if data, ok := result["data"].(map[string]interface{}); ok {
                if errorMsg, ok := data["error_message"].(string); ok {
                    fmt.Printf("📋 پیام خطا: %s\n", errorMsg)
                }
            }
        }
    }
}
```

# 📊 لاگ و مانیتورینگ

سیستم لاگ پیشرفته

```go
func setupLogging() {
    // ایجاد فایل لاگ
    logFile, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("خطا در ایجاد فایل لاگ: ", err)
    }
    
    // تنظیم خروجی لاگ
    log.SetOutput(logFile)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func logMessage(m *rubika.Message) {
    log.Printf("📨 پیام جدید - کاربر: %s, چت: %s, متن: %s", 
        m.SenderID, m.ChatID, m.Text)
}

func logError(operation string, err error) {
    log.Printf("❌ خطا در %s: %v", operation, err)
}
```

# 🎯 مثال‌های کاربردی

ربات پشتیبانی

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    bot := rubika.NewRobot("SUPPORT_BOT_TOKEN")
    
    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        switch m.Text {
        case "/start":
            sendWelcomeMessage(r, m.ChatID)
        case "📞 تماس":
            sendContactInfo(r, m.ChatID)
        case "📋 خدمات":
            sendServicesList(r, m.ChatID)
        case "🆘 پشتیبانی فوری":
            handleUrgentSupport(r, m)
        default:
            forwardToSupportTeam(r, m)
        }
    })
    
    fmt.Println("🤖 ربات پشتیبانی فعال شد")
    bot.Run()
}
```

ربات نظرسنجی

```go
package main

import (
    "fmt"
)

func main() {
    bot := rubika.NewRobot("POLL_BOT_TOKEN")
    
    bot.OnMessage(func(r *rubika.Robot, m *rubika.Message) {
        if m.Text == "/start" {
            createPoll(r, m.ChatID)
        }
    })
    
    bot.OnCallback("vote_yes", func(r *rubika.Robot, m *rubika.Message) {
        r.SendMessage(m.ChatID, "✅ نظر مثبت شما ثبت شد!", nil)
    })
    
    bot.OnCallback("vote_no", func(r *rubika.Robot, m *rubika.Message) {
        r.SendMessage(m.ChatID, "❌ نظر منفی شما ثبت شد!", nil)
    })
    
    fmt.Println("🗳️ ربات نظرسنجی فعال شد")
    bot.Run()
}
```

# 🔒 امنیت و بهترین روش‌ها

مدیریت توکن

```go
// خواندن توکن از محیطی امن
func getTokenFromEnv() string {
    token := os.Getenv("RUBIKA_BOT_TOKEN")
    if token == "" {
        log.Fatal("توکن ربات تنظیم نشده است")
    }
    return token
}

// استفاده از توکن
func main() {
    token := getTokenFromEnv()
    bot := rubika.NewRobot(token)
    // ...
}
```

اعتبارسنجی ورودی

```go
func isValidInput(input string) bool {
    // جلوگیری از حملات تزریق
    if len(input) > 1000 {
        return false
    }
    
    // جلوگیری از کاراکترهای خطرناک
    dangerousChars := []string{"<", ">", "script", "javascript:"}
    for _, char := range dangerousChars {
        if strings.Contains(input, char) {
            return false
        }
    }
    
    return true
}
```


# 📞 پشتیبانی 

منابع یادگیری

· 📚 مستندات رسمی روبیکا
· 💬 گروه پشتیبانی
· 🐙 مخزن GitHub

گزارش مشکلات

```bash
# گزارش باگ
git clone https://github.com/Daniyel-Vanguard/rubika-bot-go.git
cd rubika-bot-go
# ایجاد issue در GitHub
```

# 📜 لایسنس

این پروژه تحت لایسنس MIT منتشر شده است.

```text
MIT License

Copyright (c) 2024 [نام شما]

اجازه استفاده، کپی،، ادغام، انتشار، توزیع،
فروش و زیرمجموعه سازی این نرم‌افزار به شرط
ذکر نام صاحب اثر و این اجازه‌نامه در تمامی کپی‌ها
یا بخش‌های عمده نرم‌افزار داده می‌شود.
```

---

⭐ اگر این کتابخانه رو دوست داشتید، در GitHub ستاره بدید!

```bash
# مشارکت در توسعه
git fork https://github.com/Daniyel-Vanguard/rubika-bot-go.git
# ایجاد pull request
```

🎉 حالا شما آماده ساخت ربات‌های قدرتمند روبیکا هستید!
