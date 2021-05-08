package telegram

import (
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

func SendTelegramMessageUsingBotApi(bot *tgbotapi.BotAPI, chatId int64, dataToSend string) {
	data := utils.SplitDataInChunks(dataToSend)
	for i := 0; i < len(data); i++ {
		msg := tgbotapi.NewMessage(chatId, data[i])
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println("Error while sending message :- ", err)
		}
	}
}

func SendTelegramMessageUsingWebhook(chatID int64, dataToSend string) {
	data := utils.SplitDataInChunks(dataToSend)
	for i := 0; i < len(data); i++ {
		// Create the request body struct
		reqBody := &models.SendMessageReqBody{
			ChatID: chatID,
			Text:   data[i],
		}
		// Create the JSON body from the struct
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("Error Marshal:- ", err)
		} else {
			botToken := utils.ReadTokenFromFile()
			url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

			// Send a post request with your token
			res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				fmt.Println("Error Post:- ", err)
			}
			if res.StatusCode != http.StatusOK {
				fmt.Println("Error non 200 status code:- ", err)
			}
		}
	}
}
