package model

type troll bool

//go:generate go run github.com/MrWormHole/mapperx -source=github.com/MrWormHole/mapperx/examples/model.Admin -target=github.com/MrWormHole/mapperx/examples/model.User -output=../mapperx
type Admin struct {
	Name    string
	ID      string
	Country string
	Score   string `mapperx:"Highscore"`
	troll 	troll
	Friends []string
}
