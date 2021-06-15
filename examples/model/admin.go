package model

//go:generate mapperx github.com/MrWormHole/mapperx/examples/model.Admin github.com/MrWormHole/mapperx/examples/model.User
type Admin struct {
	Name    string
	ID      string
	Country string
	Score   string `mapperx:"Highscore"`
}
