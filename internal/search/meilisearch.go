package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MeiliClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewMeiliClient(baseURL, apiKey string) *MeiliClient {
	return &MeiliClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type MessageDocument struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type SearchResult struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

func (c *MeiliClient) do(ctx context.Context, method, path string, body any) ([]byte, error) {
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
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

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
		return nil, fmt.Errorf("meilisearch error (status %d): %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// EnsureIndex creates the messages index and configures filterable attributes.
func (c *MeiliClient) EnsureIndex(ctx context.Context) error {
	c.do(ctx, http.MethodPost, "/indexes", map[string]string{
		"uid":        "messages",
		"primaryKey": "id",
	})

	// Set filterable attributes for channel_id filtering.
	_, err := c.do(ctx, http.MethodPut, "/indexes/messages/settings/filterable-attributes",
		[]string{"channel_id"})
	return err
}

// IndexMessage adds a message document to the search index.
func (c *MeiliClient) IndexMessage(ctx context.Context, doc *MessageDocument) error {
	_, err := c.do(ctx, http.MethodPost, "/indexes/messages/documents", []any{doc})
	return err
}

// Search performs a full-text search on messages.
func (c *MeiliClient) Search(ctx context.Context, query string, channelIDs []string) ([]*SearchResult, error) {
	body := map[string]any{
		"q":     query,
		"limit": 50,
	}

	if len(channelIDs) > 0 {
		// Build filter: channel_id = 'id1' OR channel_id = 'id2'
		filters := make([]string, len(channelIDs))
		for i, id := range channelIDs {
			filters[i] = fmt.Sprintf("channel_id = '%s'", id)
		}
		body["filter"] = filters
	}

	respBody, err := c.do(ctx, http.MethodPost, "/indexes/messages/search", body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Hits []*SearchResult `json:"hits"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}
	return resp.Hits, nil
}
