package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func init() {
	// init logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// init env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file load error: %v", err)
	}
}

func main() {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/callback", func(c *fiber.Ctx) error {
		// convert *fiber.Ctx to *http.Request
		req, err := adaptor.ConvertRequest(c, false)
		if err != nil {
			logrus.Errorf("translate failed: %s", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		events, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			if err == webhook.ErrInvalidSignature {
				return c.SendStatus(fiber.StatusBadRequest)
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		for _, event := range events.Events {
			// 事件類型
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					// text message event
				}
			}
		}

		return c.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT"))))

}
