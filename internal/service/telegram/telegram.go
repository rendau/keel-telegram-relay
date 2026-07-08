package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultApiUrl = "https://api.telegram.org"

type Client struct {
	httpClient *http.Client
	apiUrl     string
	token      string
	chatId     string
}

func New(apiUrl, token, chatId string) *Client {
	if apiUrl == "" {
		apiUrl = defaultApiUrl
	}
	apiUrl = strings.TrimRight(apiUrl, "/")

	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		apiUrl:     apiUrl,
		token:      token,
		chatId:     chatId,
	}
}

type sendMessageReq struct {
	ChatId                string `json:"chat_id"`
	Text                  string `json:"text"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func (c *Client) Send(ctx context.Context, text string) error {
	reqBody, err := json.Marshal(sendMessageReq{
		ChatId:                c.chatId,
		Text:                  text,
		DisableWebPagePreview: true,
	})
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", c.apiUrl, c.token)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rep, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient.Do: %w", err)
	}
	defer rep.Body.Close()

	repBody, _ := io.ReadAll(rep.Body)

	if rep.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api status %d: %s", rep.StatusCode, string(repBody))
	}

	return nil
}
