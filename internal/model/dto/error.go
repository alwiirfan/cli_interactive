package dto

type ErrorResponse struct {
	Massage string `json:"massage"`
	Status  int    `json:"status"`
}
