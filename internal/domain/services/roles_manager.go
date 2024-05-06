package services

import "context"

type RoleManager interface {
	AssignDefaultRole(ctx context.Context, userId int64, appId int) (roleId int, err error)
}
