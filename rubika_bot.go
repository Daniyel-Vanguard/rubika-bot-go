package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Message struct {
	Bot       *Robot
	ChatID    string
	MessageID string
	SenderID  string
	Text      string
	RawData   map[string]interface{}
}

type InlineMessage struct {
	Bot     *Robot
	RawData map[string]interface{}
}

type Robot struct {
	Token              string
	Timeout            time.Duration
	Auth               string
	SessionName        string
	Key                string
	Platform           string
	OffsetID           string
	Client             *http.Client
	MessageHandler     func(*Robot, *Message)
	CallbackHandlers   []CallbackHandler
	InlineQueryHandler func(*Robot, *InlineMessage)
	WebhookURL         string
	WebhookServer      *http.Server
	mu                 sync.Mutex
	IsWebhook          bool
	PHPWebhookURL      string
}

type CallbackHandler struct {
	ButtonID string
	Handler  func(*Robot, *Message)
}

type WebhookMessage struct {
	ReceivedAt string                 `json:"received_at"`
	Data       map[string]interface{} `json:"data"`
}

const API_URL = "https://botapi.rubika.ir/v3"

func NewRobot(token string, options ...func(*Robot)) *Robot {
	robot := &Robot{
		Token:    token,
		Timeout:  10 * time.Second,
		Platform: "web",
		Client:   &http.Client{Timeout: 10 * time.Second},
	}

	for _, option := range options {
		option(robot)
	}

	return robot
}

func WithTimeout(timeout time.Duration) func(*Robot) {
	return func(r *Robot) {
		r.Timeout = timeout
		r.Client = &http.Client{Timeout: timeout}
	}
}

func WithAuth(auth string) func(*Robot) {
	return func(r *Robot) {
		r.Auth = auth
	}
}

func WithSessionName(sessionName string) func(*Robot) {
	return func(r *Robot) {
		r.SessionName = sessionName
	}
}

func WithKey(key string) func(*Robot) {
	return func(r *Robot) {
		r.Key = key
	}
}

func WithPlatform(platform string) func(*Robot) {
	return func(r *Robot) {
		r.Platform = platform
	}
}

func WithWebhook(webhookURL string) func(*Robot) {
	return func(r *Robot) {
		r.WebhookURL = webhookURL
		r.IsWebhook = true
	}
}

func WithPHPWebhook(phpWebhookURL string) func(*Robot) {
	return func(r *Robot) {
		r.PHPWebhookURL = phpWebhookURL
		r.IsWebhook = true
	}
}

func (r *Robot) post(method string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s/%s", API_URL, r.Token, method)
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Robot) OnMessage(handler func(*Robot, *Message)) {
	r.MessageHandler = handler
}

func (r *Robot) OnCallback(buttonID string, handler func(*Robot, *Message)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.CallbackHandlers = append(r.CallbackHandlers, CallbackHandler{
		ButtonID: buttonID,
		Handler:  handler,
	})
}

func (r *Robot) OnInlineQuery(handler func(*Robot, *InlineMessage)) {
	r.InlineQueryHandler = handler
}

// پردازش به‌روزرسانی‌ها
func (r *Robot) processUpdate(update map[string]interface{}) {
	updateType, ok := update["type"].(string)
	if !ok {
		return
	}

	if updateType == "ReceiveQuery" && r.InlineQueryHandler != nil {
		inlineMsg, ok := update["inline_message"].(map[string]interface{})
		if ok {
			context := &InlineMessage{
				Bot:     r,
				RawData: inlineMsg,
			}
			go r.InlineQueryHandler(r, context)
		}
		return
	}

	if updateType == "NewMessage" {
		newMessage, ok := update["new_message"].(map[string]interface{})
		if !ok {
			return
		}

		chatID, _ := update["chat_id"].(string)
		messageID, _ := newMessage["message_id"].(string)
		senderID, _ := newMessage["sender_id"].(string)
		text, _ := newMessage["text"].(string)

		if msgTime, ok := newMessage["time"].(float64); ok {
			if time.Since(time.Unix(int64(msgTime), 0)) > 20*time.Second {
				return
			}
		}

		context := &Message{
			Bot:       r,
			ChatID:    chatID,
			MessageID: messageID,
			SenderID:  senderID,
			Text:      text,
			RawData:   newMessage,
		}

		if auxData, ok := newMessage["aux_data"].(map[string]interface{}); ok {
			if buttonID, ok := auxData["button_id"].(string); ok {
				r.mu.Lock()
				handlers := make([]CallbackHandler, len(r.CallbackHandlers))
				copy(handlers, r.CallbackHandlers)
				r.mu.Unlock()
				
				for _, handler := range handlers {
					if handler.ButtonID == "" || handler.ButtonID == buttonID {
						go handler.Handler(r, context)
						return
					}
				}
			}
		}

		if r.MessageHandler != nil {
			go r.MessageHandler(r, context)
		}
	}
}

func (r *Robot) GetUpdates(offsetID string, limit int) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if offsetID != "" {
		data["offset_id"] = offsetID
	}
	if limit > 0 {
		data["limit"] = limit
	}

	return r.post("getUpdates", data)
}

func (r *Robot) Run() {
	fmt.Println("🤖 Rubika Bot started running in Polling mode...")

	if r.OffsetID == "" {
		updates, err := r.GetUpdates("", 100)
		if err == nil {
			if data, ok := updates["data"].(map[string]interface{}); ok {
				if nextOffset, ok := data["next_offset_id"].(string); ok {
					r.OffsetID = nextOffset
					fmt.Printf("📊 Offset initialized to: %s\n", r.OffsetID)
				}
			}
		}
	}

	for {
		updates, err := r.GetUpdates(r.OffsetID, 100)
		if err != nil {
			fmt.Printf("❌ Error getting updates: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if data, ok := updates["data"].(map[string]interface{}); ok {
			if updatesList, ok := data["updates"].([]interface{}); ok {
				fmt.Printf("📨 Received %d updates\n", len(updatesList))
				for _, update := range updatesList {
					if updateMap, ok := update.(map[string]interface{}); ok {
						r.processUpdate(updateMap)
					}
				}
			}

			if nextOffset, ok := data["next_offset_id"].(string); ok {
				r.OffsetID = nextOffset
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (r *Robot) StartWebhookServer(port string) error {
	if !r.IsWebhook {
		return fmt.Errorf("webhook is not configured")
	}

	_, err := r.post("updateBotEndpoints", map[string]interface{}{
		"url":  r.WebhookURL,
		"type": "Webhook",
	})
	if err != nil {
		return fmt.Errorf("failed to set webhook: %v", err)
	}

	fmt.Printf("🌐 Webhook set to: %s\n", r.WebhookURL)
	fmt.Printf("🚀 Starting webhook server on port %s\n", port)

	http.HandleFunc("/webhook", r.webhookHandler)
	
	r.WebhookServer = &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	return r.WebhookServer.ListenAndServe()
}

func (r *Robot) StartPHPWebhook() error {
	if r.PHPWebhookURL == "" {
		return fmt.Errorf("PHP webhook URL is not configured")
	}

	_, err := r.post("updateBotEndpoints", map[string]interface{}{
		"url":  r.PHPWebhookURL,
		"type": "Webhook",
	})
	if err != nil {
		return fmt.Errorf("failed to set PHP webhook: %v", err)
	}

	fmt.Printf("🌐 PHP Webhook set to: %s\n", r.PHPWebhookURL)
	fmt.Println("✅ PHP Webhook activated! Messages will be sent to your PHP script.")

	return nil
}

func (r *Robot) webhookHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var update map[string]interface{}
	if err := json.Unmarshal(body, &update); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	go r.processUpdate(update)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (r *Robot) StopWebhook() error {
	if r.WebhookServer != nil {
		return r.WebhookServer.Close()
	}
	return nil
}

func (r *Robot) SendMessage(chatID, text string, options ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	if len(options) > 0 {
		for key, value := range options[0] {
			payload[key] = value
		}
	}

	return r.post("sendMessage", payload)
}

func (r *Robot) GetMe() (map[string]interface{}, error) {
	return r.post("getMe", nil)
}

func (r *Robot) GetChat(chatID string) (map[string]interface{}, error) {
	return r.post("getChat", map[string]interface{}{
		"chat_id": chatID,
	})
}

func (r *Robot) UploadFile(filePath, mediaType string) (string, error) {
	uploadResult, err := r.post("requestSendFile", map[string]interface{}{
		"type": mediaType,
	})
	if err != nil {
		return "", err
	}

	data, ok := uploadResult["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid upload response")
	}

	uploadURL, ok := data["upload_url"].(string)
	if !ok {
		return "", fmt.Errorf("upload URL not found")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	// ارسال فایل
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := r.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	resultData, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid upload response")
	}

	fileID, ok := resultData["file_id"].(string)
	if !ok {
		return "", fmt.Errorf("file ID not found")
	}

	return fileID, nil
}

func (r *Robot) SendFile(chatID, filePath, mediaType string, options ...map[string]interface{}) (map[string]interface{}, error) {
	fileID, err := r.UploadFile(filePath, mediaType)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"chat_id": chatID,
		"file_id": fileID,
	}

	if len(options) > 0 {
		for key, value := range options[0] {
			payload[key] = value
		}
	}

	return r.post("sendFile", payload)
}

func (r *Robot) SendImage(chatID, filePath string, options ...map[string]interface{}) (map[string]interface{}, error) {
	return r.SendFile(chatID, filePath, "Image", options...)
}

func (r *Robot) SendDocument(chatID, filePath string, options ...map[string]interface{}) (map[string]interface{}, error) {
	return r.SendFile(chatID, filePath, "File", options...)
}

func (r *Robot) SendMusic(chatID, filePath string, options ...map[string]interface{}) (map[string]interface{}, error) {
	return r.SendFile(chatID, filePath, "Music", options...)
}

func (r *Robot) SendVoice(chatID, filePath string, options ...map[string]interface{}) (map[string]interface{}, error) {
	return r.SendFile(chatID, filePath, "Voice", options...)
}

func (r *Robot) SendGif(chatID, filePath string, options ...map[string]interface{}) (map[string]interface{}, error) {
	return r.SendFile(chatID, filePath, "Gif", options...)
}

func (r *Robot) DeleteMessage(chatID, messageID string) (map[string]interface{}, error) {
	return r.post("deleteMessage", map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	})
}

func (r *Robot) SendLocation(chatID string, latitude, longitude float64, options ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"chat_id":   chatID,
		"latitude":  latitude,
		"longitude": longitude,
	}

	if len(options) > 0 {
		for key, value := range options[0] {
			payload[key] = value
		}
	}

	return r.post("sendLocation", payload)
}

// ارسال مخاطب
func (r *Robot) SendContact(chatID, firstName, lastName, phoneNumber string, options ...map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"chat_id":      chatID,
		"first_name":   firstName,
		"last_name":    lastName,
		"phone_number": phoneNumber,
	}

	if len(options) > 0 {
		for key, value := range options[0] {
			payload[key] = value
		}
	}

	return r.post("sendContact", payload)
}

// ارسال نظرسنجی
func (r *Robot) SendPoll(chatID, question string, options []string) (map[string]interface{}, error) {
	return r.post("sendPoll", map[string]interface{}{
		"chat_id":  chatID,
		"question": question,
		"options":  options,
	})
}

// ویرایش پیام
func (r *Robot) EditMessageText(chatID, messageID, text string) (map[string]interface{}, error) {
	return r.post("editMessageText", map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	})
}

func (r *Robot) ForwardMessage(fromChatID, messageID, toChatID string, disableNotification bool) (map[string]interface{}, error) {
	return r.post("forwardMessage", map[string]interface{}{
		"from_chat_id":          fromChatID,
		"message_id":            messageID,
		"to_chat_id":            toChatID,
		"disable_notification":  disableNotification,
	})
}

func CreateInlineButton(text, buttonID string, buttonType ...string) map[string]interface{} {
	btnType := "Simple"
	if len(buttonType) > 0 {
		btnType = buttonType[0]
	}

	return map[string]interface{}{
		"id":          buttonID,
		"type":        btnType,
		"button_text": text,
	}
}

func CreateButtonRow(buttons ...map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"buttons": buttons,
	}
}

func CreateInlineKeypad(rows []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"rows": rows,
	}
}
