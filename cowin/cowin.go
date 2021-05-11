package cowin

import (
	"VaccineAvailabilityTelegramBot/api"
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

func FetchDataByPinCode(pincode int) string {
	url := api.GetUrlByPincode(pincode)
	return FetchDataForVaccineCenters(url)
}

func FetchDataByDistrictId(districtId int) string {
	//vadodara_corporation_dist_id := 777
	url := api.GetUrlByDistrictId(districtId)
	return FetchDataForVaccineCenters(url)
}

func FetchDataForVaccineCenters(url string) string {
	vaccineCentersDataInBytes := utils.FetchDataInBytesFromGetApiCall(url)
	if vaccineCentersDataInBytes == nil {
		fmt.Println("Unable to fetch data from api FetchDataForVaccineCenters")
		return ""
	} else {
		var vaccineCentersData models.ResponseCenters
		err := json.Unmarshal(vaccineCentersDataInBytes, &vaccineCentersData)
		if err != nil {
			fmt.Println("Error ReadData (Unmarshal) :- ", err)
			return ""
		}
		return ParseResponseCenters(vaccineCentersData)
	}
}

func ParseResponseCenters(vaccineCentersData models.ResponseCenters) string {
	messageToBeReturned := ""
	for i := 0; i < len(vaccineCentersData.Centers); i++ {
		center := vaccineCentersData.Centers[i]
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
