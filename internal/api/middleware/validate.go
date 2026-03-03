package middleware

import (
	"regexp"
)

var (
	usernameRegex    = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	channelNameRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,48}[a-z0-9]$`)
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func ValidateMessageContent(content string) bool {
	return len(content) > 0 && len(content) <= 10000
}

func ValidateChannelName(name string) bool {
	return channelNameRegex.MatchString(name)
}
