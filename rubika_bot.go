// Package rubika provides a comprehensive Go library for Rubika Bot API
// Version: 1.0.0
// Author: Daniyel Vanguard
// License: MIT
package rubika

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

const API_URL = "https://botapi.rubika.ir/v3"

type ButtonStyle string

const (
	Primary   ButtonStyle = "Primary"
	Secondary ButtonStyle = "Secondary"
	Danger    ButtonStyle = "Danger"
	Success   ButtonStyle = "Success"
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

func CreateInlineButton(text, buttonID string, style ButtonStyle) map[string]interface{} {
	return map[string]interface{}{
		"text":      text,
		"button_id": buttonID,
		"style":     string(style),
	}
}

func CreateButtonRow(buttons ...map[string]interface{}) []map[string]interface{} {
	return buttons
}

func CreateInlineKeypad(rows []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"rows": rows,
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
