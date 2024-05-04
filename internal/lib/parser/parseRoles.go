package parser

import "auth_service/internal/domain/models"

type JwtRoles struct {
	RolesName []string
}

func ParseUserRoles(userRoles models.UserRoles) JwtRoles {
	var jwtRoles JwtRoles
	for _, role := range userRoles.Roles {
		jwtRoles.RolesName = append(jwtRoles.RolesName, role.Name)
	}
	return jwtRoles
}
