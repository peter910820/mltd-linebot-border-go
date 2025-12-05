package handlers

import (
	"mltd-linebot-border-go/matsurihime"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/sirupsen/logrus"
)

func TextMessageHandler(bot *messaging_api.MessagingApiAPI, e webhook.MessageEvent, message webhook.TextMessageContent) {
	_, err := matsurihime.GetEvents()
	if err != nil {
		postError(bot, e, err)
		return
	}
}

func postError(bot *messaging_api.MessagingApiAPI, e webhook.MessageEvent, err error) {
	logrus.Error(err)
	_, err = bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: "目前發生錯誤，請稍後再試",
				},
			},
		},
	)

	if err != nil {
		logrus.Error(err)
	}
}
