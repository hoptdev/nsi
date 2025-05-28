package models

type GrantType string

const (
	ReadOnly GrantType = "read"
	Update   GrantType = "update"
	Admin    GrantType = "admin"
)

type AccessRight struct {
	Id          int
	UserId      *int
	UserGroupId *int
	AccessToken *string
	Type        GrantType
}
