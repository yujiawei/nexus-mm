package wkim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL      string
	managerToken string
	httpClient   *http.Client
}

func NewClient(baseURL, managerToken string) *Client {
	return &Client{
		baseURL:      baseURL,
		managerToken: managerToken,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", c.managerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("wukongim api error (status %d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// RegisterUser registers a user in WuKongIM.
func (c *Client) RegisterUser(ctx context.Context, uid, token string) error {
	body := map[string]string{"uid": uid, "token": token}
	_, err := c.do(ctx, http.MethodPost, "/user/token_update", body)
	return err
}

// CreateChannel creates a group channel in WuKongIM.
func (c *Client) CreateChannel(ctx context.Context, channelID string, channelType int, subscribers []string) error {
	body := map[string]any{
		"channel_id":  channelID,
		"channel_type": channelType,
		"subscribers":  subscribers,
	}
	_, err := c.do(ctx, http.MethodPost, "/channel", body)
	return err
}

// AddSubscribers adds members to a WuKongIM channel.
func (c *Client) AddSubscribers(ctx context.Context, channelID string, channelType int, uids []string) error {
	body := map[string]any{
		"channel_id":   channelID,
		"channel_type": channelType,
		"subscribers":  uids,
	}
	_, err := c.do(ctx, http.MethodPost, "/channel/subscriber_add", body)
	return err
}

// RemoveSubscribers removes members from a WuKongIM channel.
func (c *Client) RemoveSubscribers(ctx context.Context, channelID string, channelType int, uids []string) error {
	body := map[string]any{
		"channel_id":   channelID,
		"channel_type": channelType,
		"subscribers":  uids,
	}
	_, err := c.do(ctx, http.MethodPost, "/channel/subscriber_remove", body)
	return err
}

// SendMessage sends a message through WuKongIM.
type SendMsgReq struct {
	FromUID     string `json:"from_uid"`
	ChannelID   string `json:"channel_id"`
	ChannelType int    `json:"channel_type"`
	Payload     []byte `json:"payload"`
}

func (c *Client) SendMessage(ctx context.Context, req *SendMsgReq) error {
	_, err := c.do(ctx, http.MethodPost, "/message/send", req)
	return err
}
