package dto

type CharactersResponse struct {
	LastSync   string `json:"lastSync"`
	Characters map[string][]string `json:"characters"`
}//@Name CharactersResponse
