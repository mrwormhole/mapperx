package model

type User struct {
	Name             string
	ID               string
	Country          string
	Highscore        string
	troll            troll
	Friends          []string
	GivenPermissions map[string]string
}
