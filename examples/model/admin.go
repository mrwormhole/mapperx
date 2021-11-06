package model

type troll bool

//go:generate go run github.com/MrWormHole/mapperx -source=github.com/MrWormHole/mapperx/examples/model.Admin -target=github.com/MrWormHole/mapperx/examples/model.User -directory=../mapperx -filename=XXX_ADMIN_TO_USER_XXX
// OR
///go:generate go run github.com/MrWormHole/mapperx -source=github.com/MrWormHole/mapperx/examples/model.Admin -target=github.com/MrWormHole/mapperx/examples/model.User
type Admin struct {
	Name        string
	ID          string
	Country     string
	Score       string `mapperx:"Highscore"`
	troll       troll
	Friends     []string
	Permissions map[string]string `mapperx:"GivenPermissions"`
	Minion *Minion
}
