package models

type ResponseTunnels struct {
	Tunnel []Tunnel `json:"tunnels"`
}

type Tunnel struct {
	PublicUrl string `json:"public_url"`
}
