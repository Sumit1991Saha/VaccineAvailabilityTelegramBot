package cowin

import (
	"VaccineAvailabilityTelegramBot/api"
	"VaccineAvailabilityTelegramBot/models"
	"VaccineAvailabilityTelegramBot/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
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

func FetchStateDataAlongWithDistrictData() {
	url := api.GetUrlForAllStateData()
	statesDataInBytes := utils.FetchDataInBytesFromGetApiCall(url)
	if statesDataInBytes == nil {
		fmt.Println("Missing States information")
		return
	}
	var statesData models.ResponseStates
	err := json.Unmarshal(statesDataInBytes, &statesData)
	if err != nil {
		fmt.Println("Error ReadData (Unmarshal) for state data:- ", err)
		return
	}
	WriteStateData(statesData.States)
	for _, state := range statesData.States {
		stateName := state.StateName
		stateId := state.StateId
		urlToFetchAllDistricts := api.GetUrlByStateIdForDistrictData(stateId)
		districtDataInBytes := utils.FetchDataInBytesFromGetApiCall(urlToFetchAllDistricts)
		if districtDataInBytes == nil {
			fmt.Println("Missing district data for ", stateName)
			return
		}

		var districtsData models.ResponseDistricts
		err := json.Unmarshal(districtDataInBytes, &districtsData)
		if err != nil {
			fmt.Println("Error ReadData (Unmarshal) for district data :- ", err)
		} else {
			WriteDistrictData(stateName, districtsData.Districts)
		}
	}
}

func WriteStateData(states []models.State) {
	f, e := os.Create("./States.csv")
	if e != nil {
		fmt.Println(e)
	} else {
		writer := csv.NewWriter(f)
		var stateData = [][]string{
			{"State-Id", "State-Name"},
		}

		for _, state := range states {
			stateName := state.StateName
			stateId := state.StateId
			stateData = append(stateData, []string{strconv.Itoa(stateId), stateName})
		}

		e = writer.WriteAll(stateData)
		if e != nil {
			fmt.Println(e)
		}
	}
}

func WriteDistrictData(stateName string, districts []models.District) {
	f, e := os.Create(fmt.Sprintf("./%s.csv", stateName))
	if e != nil {
		fmt.Println(e)
	} else {
		writer := csv.NewWriter(f)
		var stateData = [][]string{
			{"District-Id", "Distrcit-Name"},
		}

		for _, district := range districts {
			districtName := district.DistrictName
			districtId := district.DistrictId
			stateData = append(stateData, []string{strconv.Itoa(districtId), districtName})
		}

		e = writer.WriteAll(stateData)
		if e != nil {
			fmt.Println(e)
		}
	}
}
