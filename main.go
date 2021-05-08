package main

import (
	"VaccineAvailabilityTelegramBot/cowin"
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/telegram"
	"VaccineAvailabilityTelegramBot/utils"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strconv"
)

func main() {
	//usingGetUpdates()
	usingWebhook()
}

func usingWebhook() {
	SetWebhook()
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}

func SetWebhook() {
	//ngrokTunnelUrl := "http://localhost:4040/api/tunnels"
	//http.Get(ngrokTunnelUrl)
}

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &models.WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	if val, err := strconv.Atoi(body.Message.Text); err == nil {
		dataToSend := ""
		/*if len(update.Message.Text) != 6 {
			dataToSend = "PLease enter a valid pincode"
		} else {
			dataToSend = FetchDataByPinCode(val)
		}*/
		dataToSend = cowin.FetchDataByDistrictId(val)
		if dataToSend == "" {
			fmt.Println("Error empty string, so not sending data")
		} else {
			fmt.Println(dataToSend)
			telegram.SendTelegramMessageUsingWebhook(body.Message.Chat.ID, dataToSend)
		}
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

func usingGetUpdates() {
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
			/*if len(update.Message.Text) != 6 {
				dataToSend = "PLease enter a valid pincode"
			} else {
				dataToSend = FetchDataByPinCode(val)
			}*/
			dataToSend = cowin.FetchDataByDistrictId(val)
			if dataToSend == "" {
				fmt.Println("Error empty string, so not sending data")
			} else {
				telegram.SendTelegramMessageUsingBotApi(bot, update.Message.Chat.ID, dataToSend)
			}
		}
	}
}




