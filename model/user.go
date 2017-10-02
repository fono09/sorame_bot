package model

type User struct {
	ID             int64  `json: "id"`
	ScreenName     string `json: "screen_name"`
	ConsumerKey    string
	ConsumerSecret string
	Sorames        []Sorame `json: "sorames"`
}
