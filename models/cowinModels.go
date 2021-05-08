package models

type ResponseCenters struct {
	Centers []Center `json:"centers"`
}

type Center struct {
	CenterId     int       `json:"center_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	StateName    string    `json:"state_name"`
	DistrictName string    `json:"district_name"`
	BlockName    string    `json:"block_name"`
	Pincode      int       `json:"pincode"`
	Lat          int       `json:"lat"`
	Long         int       `json:"long"`
	From         string    `json:"from"`
	To           string    `json:"to"`
	FeeType      string    `json:"fee_type"`
	Sessions     []Session `json:"sessions"`
}

type Session struct {
	SessionId         string   `json:"session_id"`
	Date              string   `json:"date"`
	AvailableCapacity int      `json:"available_capacity"`
	MinAgeLimit       int      `json:"min_age_limit"`
	Vaccine           string   `json:"vaccine"`
	Slots             []string `json:"slots"`
}

type ResponseStates struct {
	States []State `json:"states"`
}

type State struct {
	StateId   string `json:"state_id"`
	StateName string `json:"state_name"`
}

type ResponseDistricts struct {
	District []District `json:"districts"`
}

type District struct {
	DistrictId   string `json:"district_id"`
	DistrictName string `json:"district_name"`
}
