package roles

import (
	"auth_service/internal/domain/models"
	"context"
	"fmt"
	"log/slog"
)

type RoleManager struct {
	log           *slog.Logger
	rolesProvider RoleHandler
}

type RoleHandler interface {
	GetRoleByApp(ctx context.Context, appId int64, name string) (models.Role, error)
	SaveRole(ctx context.Context, userId int64, roleId int) error
}

func New(logger *slog.Logger, rolesProvider RoleHandler) *RoleManager {
	return &RoleManager{
		log:           logger,
		rolesProvider: rolesProvider,
	}
}

func (r *RoleManager) AssignDefaultRole(ctx context.Context, userId int64, appId int) (roleId int, err error) {
	const op = "roleManager.AssignDefaultRole"
	//check is app allowNewRole

	//check default role

	//now instead search for user

	role, err := r.rolesProvider.GetRoleByApp(ctx, int64(appId), "user")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = r.rolesProvider.SaveRole(ctx, userId, role.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return role.ID, nil

}
