package models

type GrantType string

const (
	ReadOnly GrantType = "read"
	Update   GrantType = "update"
	Admin    GrantType = "admin"
)

var ranks = map[GrantType]int{ReadOnly: 0, Update: 1, Admin: 2} //make(map[string]int)

func (t GrantType) ToInt() int {
	return ranks[t]
}

type AccessRight struct {
	Id          int
	UserId      *int
	UserGroupId *int
	AccessToken *string
	Type        GrantType
}
