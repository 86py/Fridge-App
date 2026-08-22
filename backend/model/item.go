package model

type Item struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Quantity       int    `json:"quantity"`
	Category       string `json:"category"`
	ExpirationDate string `json:"expiration_date"`
	CreatedAt      string `json:"created_at"`
}
