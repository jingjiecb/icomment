package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"claws.top/icomment/model"
)

const (
	BarkAPIURL         = "https://api.day.app"
	MaxContentPreview  = 100 // 通知中显示的最大字符数
	NotificationTimeout = 5 * time.Second
)

type BarkNotifier struct {
	deviceKey string
	client    *http.Client
}

func NewBarkNotifier(deviceKey string) *BarkNotifier {
	if deviceKey == "" {
		return nil
	}
	return &BarkNotifier{
		deviceKey: deviceKey,
		client: &http.Client{
			Timeout: NotificationTimeout,
		},
	}
}

type BarkMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
	Group string `json:"group,omitempty"`
}

func (b *BarkNotifier) SendNewCommentNotification(comment *model.Comment) error {
	if b == nil {
		return nil // Bark not configured
	}

	// 截取评论内容预览
	contentPreview := comment.Content
	if len(contentPreview) > MaxContentPreview {
		// 按字符截取，避免截断 emoji
		runes := []rune(contentPreview)
		if len(runes) > MaxContentPreview {
			contentPreview = string(runes[:MaxContentPreview]) + "..."
		}
	}

	// 构建通知消息
	title := "新评论"
	if comment.ParentID != nil {
		title = "新回复"
	}

	body := fmt.Sprintf("来自 %s:\n%s", comment.Nickname, contentPreview)

	message := BarkMessage{
		Title: title,
		Body:  body,
		URL:   comment.ArticleURL,
		Group: "iComment",
	}

	// 发送通知
	return b.send(message)
}

func (b *BarkNotifier) send(message BarkMessage) error {
	url := fmt.Sprintf("%s/%s", BarkAPIURL, b.deviceKey)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal bark message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create bark request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send bark notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark API returned status %d", resp.StatusCode)
	}

	return nil
}
