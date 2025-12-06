package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/sirupsen/logrus"

	"mltd-linebot-border-go/common"
	"mltd-linebot-border-go/matsurihime"
)

var (
	timeFormat = "2006-01-02 15:04:05"

	reEvent     = regexp.MustCompile(`(?i)^event-[a-z]{2}$`)
	reEventRank = regexp.MustCompile(`(?i)^event-[a-z]{2}-[0-9]+$`)

	defaultFormat = map[string]string{
		"eventPoint":  "1,2,3,100,2500,5000,10000,25000,50000,100000",
		"highScore":   "1,2,3,100,2000,5000,10000,20000,100000",
		"loungePoint": "1,2,3,10,100,250,500,1000",
	}
)

func TextMessageHandler(bot *messaging_api.MessagingApiAPI, e webhook.MessageEvent, message webhook.TextMessageContent) {
	// rank
	if reEventRank.MatchString(message.Text) {
		returnMsg, err := mltdBorder(message, false)
		if err != nil {
			postError(bot, e, err)
			return
		}
		postMessage(bot, e, returnMsg)
	} else if reEvent.MatchString(message.Text) { // default
		returnMsg, err := mltdBorder(message, true)
		if err != nil {
			postError(bot, e, err)
			return
		}
		postMessage(bot, e, returnMsg)
	} else if strings.Contains(strings.ToLower(message.Text), "help") {
		returnMsg := "詳細指令請參考https://github.com/peter910820/mltd-linebot-border-go"
		postMessage(bot, e, returnMsg)
	}

}

func mltdBorder(message webhook.TextMessageContent, defaultMark bool) (string, error) {
	events, err := matsurihime.GetEvents()
	if err != nil {
		return "", err
	}

	var logType string
	var rankFormat string

	runes := []rune(message.Text)
	switch string(runes[6:8]) {
	case "pt":
		logType = "eventPoint"
		rankFormat = defaultFormat["eventPoint"]
	case "hs":
		logType = "highScore"
		rankFormat = defaultFormat["highScore"]
	case "lp":
		logType = "loungePoint"
		rankFormat = defaultFormat["loungePoint"]
	default:
		return "", common.ErrMLTDLogTypeAbnormal
	}

	if !defaultMark {
		rankFormat = string(runes[9:])
	}

	rankings, err := matsurihime.GetRankings(events[len(events)-1].ID, logType, rankFormat)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("%s\n開始時間: %s\n結束時間: %s\n(名次/分數/半小時增加量)\n\n",
		events[len(events)-1].Name,
		events[len(events)-1].Schedule.BeginAt.Format(timeFormat),
		events[len(events)-1].Schedule.EndAt.Format(timeFormat))

	rankingsData := make([]string, 0, len(rankings))
	for _, ranking := range rankings {
		dataLen := len(ranking.Data)
		prefixEmoji := "🔴"
		switch ranking.Rank {
		case 1:
			prefixEmoji = "🥇"
		case 2:
			prefixEmoji = "🥈"
		case 3:
			prefixEmoji = "🥉"
		}
		if dataLen > 1 {
			rankingsData = append(rankingsData, fmt.Sprintf("%s%d位:  %dpt(+%d)",
				prefixEmoji,
				ranking.Rank,
				ranking.Data[dataLen-1].Score,
				ranking.Data[dataLen-1].Score-ranking.Data[dataLen-2].Score))
		} else if dataLen > 0 {
			rankingsData = append(rankingsData, fmt.Sprintf("%s%d位:  %d",
				prefixEmoji,
				ranking.Rank,
				ranking.Data[dataLen-1].Score))
		}
	}
	output += strings.Join(rankingsData, "\n")
	return output, nil
}

func postMessage(bot *messaging_api.MessagingApiAPI, e webhook.MessageEvent, msg string) {
	_, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: msg,
				},
			},
		},
	)

	if err != nil {
		postError(bot, e, err)
	}
}

func postError(bot *messaging_api.MessagingApiAPI, e webhook.MessageEvent, err error) {
	logrus.Error(err)
	replyErrMsg := "目前發生錯誤，請稍後再試"
	if errors.Is(err, common.ErrMLTDLogTypeAbnormal) {
		replyErrMsg = "參數錯誤，目前僅支援pt hs lp三種"
	}
	_, err = bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: replyErrMsg,
				},
			},
		},
	)

	if err != nil {
		logrus.Error(err)
	}
}
