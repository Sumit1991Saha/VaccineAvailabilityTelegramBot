package cowin

import (
	"VaccineAvailabilityTelegramBot/api"
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func FetchDataByPinCode(pincode int) string {
	url := api.GetUrlByPincode(pincode)
	return FetchData(url)
}

func FetchDataByDistrictId(districtId int) string {
	//vadodara_corporation_dist_id := 777
	url := api.GetUrlByDistrictId(districtId)
	return FetchData(url)
}

func FetchData(url string) string {
	response, err := utils.ApiCall(url, "GET")
	if err != nil || response == nil {
		fmt.Println("Error (FetchData) :- ", err)
		return ""
	}
	vacancies := ReadDataForCenters(response)
	return vacancies
	//call("http://pokeapi.co/api/v2/pokedex/kanto/", "GET")
}

func ReadDataForCenters(response *http.Response) string {

	if response.StatusCode != http.StatusOK {
		fmt.Println("Error calling api :- " , response.Status)
		return ""
	}

	var jsonData models.ResponseCenters

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error ReadData (ReadAll) :- ", err)
		return ""
	}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		fmt.Println("Error ReadData (Unmarshal) :- ", err)
		return ""
	}

	return ParseResponseCenters(jsonData)
}

func ParseResponseCenters(jsonData models.ResponseCenters) string {
	messageToBeReturned := ""
	for i := 0; i < len(jsonData.Centers); i++ {
		center := jsonData.Centers[i]
		sessions := center.Sessions
		for j := 0; j < len(sessions); j++ {
			session := sessions[j]
			noOfDosesAvailable := session.AvailableCapacity
			if noOfDosesAvailable > 0 {
				msg := "Vaccine available on :- " + session.Date + "\n" +
					"AvailableCapacity :- " + strconv.Itoa(noOfDosesAvailable) + "\n" +
					"At center :- " + center.Name + "\n" +
					"Address :- " + center.Address + "\n"
				if session.MinAgeLimit >= 45 {
					msg = "For 45+ years, " + "\n" + msg
				} else {
					msg = "For 18+ years, " + "\n" + msg
				}
				messageToBeReturned = messageToBeReturned + msg + "\n"
			}
		}
	}
	if messageToBeReturned == "" {
		messageToBeReturned = "No slots available, PLease try again tomorrow"
	}
	return messageToBeReturned
}
