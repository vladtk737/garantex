package entity

type Trade struct {
	ID        int    `json:"id"`
	Price     string `json:"price"`
	Volume    string `json:"volume"`
	Funds     string `json:"funds"`
	Market    string `json:"market"`
	CreatedAt string `json:"created_at"`
}
