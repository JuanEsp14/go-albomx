package dto

type CollaboratorsResponse struct {
	LastSync  string   `json:"lastSync"`
	Editors   []string `json:"editors"`
	Writers   []string `json:"writers"`
	Colorists []string `json:"colorists"`
} //@Name CollaboratorsResponse
