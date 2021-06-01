package model

//go:generate go run github.com/mrwormhole/mapperx github.com/mrwormhole/mapperx/examples/model.Admin github.com/mrwormhole/mapperx/examples/model.User
type Admin struct {
	Name        string
	ID          string
	Country     string
	Score       string
	Permissions []string
}
