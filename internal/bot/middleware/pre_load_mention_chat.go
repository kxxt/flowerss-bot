package middleware

import (
	"strings"

	"github.com/indes/flowerss-bot/internal/bot/chat"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/log"

	tb "gopkg.in/telebot.v3"
)

func PreLoadMentionChat() tb.MiddlewareFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			fields := strings.Fields(c.Text())
			if len(fields) > 1 {
				chat, err := chat.GetChatByIdOrUsername(c.Bot(), fields[1])
				if err != nil {
					log.Errorf("pre load mention %s chat failed, %v", fields[1], err)
					return next(c)
				}
				c.Set(session.StoreKeyMentionChat.String(), chat)
			}
			return next(c)
		}
	}
}
