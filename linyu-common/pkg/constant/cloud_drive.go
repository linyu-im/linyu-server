package constant

// ------------------space-file相关------------------
type spaceType struct {
	User  string
	Group string
	//Org   string
}

var SpaceType = spaceType{
	User:  "user",
	Group: "group",
	//Org:   "org",
}

func (c spaceType) Validate(v string) bool {
	switch v {
	case c.User, c.Group:
		return true
	default:
		return false
	}
}
