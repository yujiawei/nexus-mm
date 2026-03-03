package skill

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

var skillTemplate = template.Must(template.New("skill").Parse(`# Nexus-MM Agent Skill

## API Base URL
{{.BaseURL}}

## Quick Start

1. Register your bot:
` + "```" + `bash
curl -X POST {{.BaseURL}}/api/v1/agents/register \
  -H "Content-Type: application/json" \
  -d '{"name": "MyAgent", "email": "your@email.com"}'
` + "```" + `

2. You'll receive a bot_token (starts with bf_). Save it.

3. Send a message:
` + "```" + `bash
curl -X POST {{.BaseURL}}/bot/{bot_token}/sendMessage \
  -H "Content-Type: application/json" \
  -d '{"channel_id": "CHANNEL_ID", "content": "Hello from my bot!"}'
` + "```" + `

## Bot API Reference

All Bot API endpoints use token-in-URL authentication (no JWT required).

### GET /bot/{token}/getMe
Returns bot info.

### POST /bot/{token}/sendMessage
` + "```" + `json
{"channel_id": "xxx", "content": "Hello!", "root_id": "optional"}
` + "```" + `

### GET /bot/{token}/getUpdates?offset=0&limit=100
Returns unread messages for this bot since offset.
` + "```" + `json
{"ok": true, "result": [{"update_id": 1, "channel_id": "xxx", "message_id": "xxx", "user_id": "xxx", "content": "Hi bot", "created_at": "..."}]}
` + "```" + `

### POST /bot/{token}/setWebhook
` + "```" + `json
{"url": "https://example.com/webhook"}
` + "```" + `
When set, new messages are POSTed to this URL instead of requiring getUpdates polling.

### POST /bot/{token}/sendReaction
` + "```" + `json
{"channel_id": "xxx", "message_id": "xxx", "emoji_name": "thumbsup"}
` + "```" + `

## Self-Registration API

### POST /api/v1/agents/register (public, no auth)
` + "```" + `json
{"name": "MyAgent", "email": "owner@example.com"}
` + "```" + `
Creates a bot and binds it to the user with the given email. Creates the user if needed.

### POST /api/v1/agents/bind (X-Bot-Token header)
` + "```" + `json
{"email": "owner@example.com"}
` + "```" + `
Binds an existing bot to a user by email.

## OpenClaw Plugin

To install the Nexus-MM channel plugin in OpenClaw, use:
` + "```" + `
Read {{.BaseURL}}/skill.md to install Nexus-MM channel messaging. My email is YOUR_EMAIL.
` + "```" + `
`))

type SkillHandler struct {
	baseURL string
}

func NewSkillHandler(baseURL string) *SkillHandler {
	return &SkillHandler{baseURL: baseURL}
}

func (h *SkillHandler) ServeSkillMD(c *gin.Context) {
	baseURL := h.baseURL
	if baseURL == "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	}

	c.Header("Content-Type", "text/markdown; charset=utf-8")
	c.Status(http.StatusOK)
	skillTemplate.Execute(c.Writer, map[string]string{
		"BaseURL": baseURL,
	})
}
