package telegram

import (
	"VaccineAvailabilityTelegramBot/api"
	"VaccineAvailabilityTelegramBot/cowin"
	"VaccineAvailabilityTelegramBot/telegram/sendMessage"
	"VaccineAvailabilityTelegramBot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
)

func StartServiceUsingGetUpdates() {
	deleteWebhookIfAny()
	botToken := utils.ReadTokenFromFile()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Error (usingGetUpdates) :- ", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if val, err := strconv.Atoi(update.Message.Text); err == nil {
			dataToSend := ""
			if len(update.Message.Text) == 6 {
				dataToSend = cowin.FetchDataByPinCode(val)
			} else if len(update.Message.Text) <=3 {
				dataToSend = cowin.FetchDataByDistrictId(val)
			} else {
				sendMessage.SendTelegramUsingBotApi(bot, update.Message.Chat.ID,
					"Please enter valid pincode or district id")
				return
			}

			if dataToSend == "" {
				fmt.Println("Error empty string, so not sending data")
				sendMessage.SendTelegramUsingBotApi(bot, update.Message.Chat.ID,
					"Unable to fetch data, please try again after sometime")
			} else {
				sendMessage.SendTelegramUsingBotApi(bot, update.Message.Chat.ID, dataToSend)
				fmt.Println("Data sent")
			}
		} else {
			sendMessage.SendTelegramUsingBotApi(bot, update.Message.Chat.ID,
				"Please enter either valid pincode or district id")
		}
	}
}

func deleteWebhookIfAny() {
	webhookToDelete := api.GetUrlToDeleteWebhook()
	response, err := http.Get(webhookToDelete)
	if err != nil {
		return
	}
	if response.StatusCode == http.StatusOK {
		fmt.Println("Webhook is deleted")
	} else {
		fmt.Println("unable to delete Webhook")
	}
}
