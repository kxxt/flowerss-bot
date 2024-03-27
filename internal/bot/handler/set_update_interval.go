package handler

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	tb "gopkg.in/telebot.v3"

	"github.com/indes/flowerss-bot/internal/bot/message"
	"github.com/indes/flowerss-bot/internal/bot/session"
	"github.com/indes/flowerss-bot/internal/core"
	"github.com/indes/flowerss-bot/internal/log"
)

type SetUpdateInterval struct {
	core *core.Core
}

func NewSetUpdateInterval(core *core.Core) *SetUpdateInterval {
	return &SetUpdateInterval{core: core}
}

func (s *SetUpdateInterval) Command() string {
	return "/setinterval"
}

func (s *SetUpdateInterval) Description() string {
	return "设置订阅刷新频率"
}

func (s *SetUpdateInterval) getMessageWithoutMention(ctx tb.Context) string {
	mention := message.MentionFromMessage(ctx.Message())
	if mention == "" {
		return ctx.Message().Payload
	}
	return strings.Replace(ctx.Message().Payload, mention, "", -1)
}

func (s *SetUpdateInterval) Handle(ctx tb.Context) error {
	msg := s.getMessageWithoutMention(ctx)
	args := strings.Split(strings.TrimSpace(msg), " ")
	if len(args) < 2 {
		return ctx.Reply("/setinterval [chatID] [interval] [sourceID] 设置订阅刷新频率（可设置多个sub id，以空格分割）(0 表示使用当前 Chat)")
	}

	var subscribeUserID int64 = 0
	var err error
	if strings.HasPrefix(args[0], "@") {
		var chatID *tb.Chat
		chatID, err = ctx.Bot().ChatByUsername(args[0])
		if err == nil {
			subscribeUserID = chatID.ID
		}
	} else {
		subscribeUserID, err = strconv.ParseInt(args[0], 10, 64)
	}
	if err != nil {
		return ctx.Reply("请输入正确的 Chat ID")
	}

	interval, err := strconv.Atoi(args[0])
	if interval <= 0 || err != nil {
		return ctx.Reply("请输入正确的抓取频率")
	}

	if subscribeUserID == 0 {
		subscribeUserID = ctx.Message().Chat.ID
		mentionChat, _ := session.GetMentionChatFromCtxStore(ctx)
		if mentionChat != nil {
			subscribeUserID = mentionChat.ID
		}
	}

	for _, id := range args[2:] {
		sourceID := cast.ToUint(id)
		if err := s.core.SetSubscriptionInterval(
			context.Background(), subscribeUserID, sourceID, interval,
		); err != nil {
			log.Errorf("SetSubscriptionInterval failed, %v", err)
			return ctx.Reply("抓取频率设置失败!")
		}
	}
	return ctx.Reply("抓取频率设置成功!")
}

func (s *SetUpdateInterval) Middlewares() []tb.MiddlewareFunc {
	return nil
}
