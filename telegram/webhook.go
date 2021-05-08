package telegram

import (
	"VaccineAvailabilityTelegramBot/api"
	"VaccineAvailabilityTelegramBot/cowin"
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/telegram/sendMessage"
	"VaccineAvailabilityTelegramBot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func StartServiceUsingWebhook() {
	SetWebhook()
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}

func SetWebhook() {
	ngrokTunnelUrl := api.AllNgrokTunnels
	response, err := utils.ApiCall(ngrokTunnelUrl, "GET")
	if err != nil || response == nil {
		fmt.Println("Error (FetchData) :- ", err)
		return
	}
	ngrokServerUrl := utils.GetPublicUrlNgrok(response)
	if ngrokServerUrl == "" {
		fmt.Println("Error Unable to set Webhook")
	} else {
		reqBody := &models.Webhook{
			Url : ngrokServerUrl,
		}
		// Create the JSON body from the struct
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("Error Marshal:- ", err)
		} else {
			url := api.GetUrlToSetWebhook()
			// Send a post request with your token
			res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				fmt.Println("Error Post:- ", err)
			}
			if res.StatusCode != http.StatusOK {
				fmt.Println("Error non 200 status code:- ", err)
			}
			if err == nil && res.StatusCode == http.StatusOK {
				fmt.Println("Webhook is set")
			}
		}

	}
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
			sendMessage.SendTelegramUsingWebhook(body.Message.Chat.ID, dataToSend)
		}
	}
}
