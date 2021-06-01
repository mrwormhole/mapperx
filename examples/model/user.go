package model

type User struct {
	Name      string
	ID        string
	Country   string `mapperx:"Score"`
	Highscore string
}
