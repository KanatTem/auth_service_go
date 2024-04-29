package models

const (
	UserRoleId    = 1
	UserRoleText  = "user"
	AdminRoleId   = 2
	AdminRoleText = "admin"
)

type Role struct {
	ID   int
	Name string
}
