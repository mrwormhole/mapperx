package model

//go:generate go run github.com/MrWormHole/mapperx github.com/MrWormHole/mapperx/examples/model.Admin github.com/MrWormHole/mapperx/examples/model.User
type Admin struct {
	Name    string
	ID      string
	Country string
	Score   string
}
