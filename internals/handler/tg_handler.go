package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	telegramAPI = "https://api.telegram.org/bot"
)

type TelegramBot struct {
	token     string
	channelID string
	client    *http.Client
}

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description,omitempty"`
}

func NewTelegramBot(token, channelID string) *TelegramBot {
	return &TelegramBot{
		token:     token,
		channelID: channelID,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (bot *TelegramBot) SendMessage(message string) error {
	endpoint := fmt.Sprintf("%s%s/sendMessage", telegramAPI, bot.token)

	data := url.Values{}
	data.Add("chat_id", bot.channelID)
	data.Add("text", message)
	data.Add("parse_mode", "HTML")

	resp, err := bot.client.PostForm(endpoint, data)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var response TelegramResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if !response.Ok {
		return fmt.Errorf("telegram API error: %s", response.Description)
	}

	return nil
}
