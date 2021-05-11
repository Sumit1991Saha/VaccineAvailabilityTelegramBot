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
	ngrokTunnelDataInBytes := utils.FetchDataInBytesFromGetApiCall(ngrokTunnelUrl)

	var ngrokTunnelData models.ResponseTunnels

	err := json.Unmarshal(ngrokTunnelDataInBytes, &ngrokTunnelData)
	if err != nil {
		fmt.Println("Error ReadData (Unmarshal) :- ", err)
	} else {
		ngrokServerUrl := ngrokTunnelData.Tunnel[0].PublicUrl
		if ngrokServerUrl == "" {
			fmt.Println("Error Unable to set Webhook")
		} else {
			reqBody := &models.Webhook{
				Url: ngrokServerUrl,
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
					fmt.Println("Error setting Webhook:- ", err)
				}
				if res.StatusCode != http.StatusOK {
					fmt.Println("Unable to set Webhook")
				}
				if err == nil && res.StatusCode == http.StatusOK {
					fmt.Println("Webhook is set")
				}
			}
		}
	}
}

// Handler This handler is called everytime telegram sends us a webhook event
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &models.TelegramRequest{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	fmt.Println("Received text :- " + body.Message.Text + " , from :- " + body.Message.From.FirstName)

	if val, err := strconv.Atoi(body.Message.Text); err == nil {
		dataToSend := ""
		if len(body.Message.Text) == 6 {
			dataToSend = cowin.FetchDataByPinCode(val)
		} else if len(body.Message.Text) <= 3 {
			dataToSend = cowin.FetchDataByDistrictId(val)
		} else {
			sendMessage.SendTelegramUsingWebhook(body.Message.Chat.Id,
				"Please enter valid pincode or district id")
			return
		}
		if dataToSend == "" {
			fmt.Println("Error empty string, so not sending data")
			sendMessage.SendTelegramUsingWebhook(body.Message.Chat.Id,
				"Unable to fetch data, please try again after sometime")
		} else {
			//fmt.Println(dataToSend)
			sendMessage.SendTelegramUsingWebhook(body.Message.Chat.Id, dataToSend)
			fmt.Println("Data sent")
		}
	} else {
		sendMessage.SendTelegramUsingWebhook(body.Message.Chat.Id,
			"Please enter either valid pincode or district id")
	}
}
