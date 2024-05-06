package postgress

import (
	"auth_service/internal/domain/models"
	"auth_service/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("postgres", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() {
	const op = "storage.postgres.Stop"

	fmt.Printf("stopping bd connection")

	err := s.db.Close()
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
	}

}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	// 1) Prepare the INSERT with RETURNING id
	query := `
        INSERT INTO users (email, pass_hash)
        VALUES ($1, $2)
        RETURNING id
    `

	// 2) Use QueryRowContext to execute and scan the new id
	var id int64
	err := s.db.QueryRowContext(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		// handle unique‐constraint violation, etc.
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, storage.ErrUserExists
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (models.User, error) {

	const op = "storage.postgres.GetUser"

	stmt, err := s.db.Prepare("select * from users where email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User

	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}

func (s *Storage) GetApp(ctx context.Context, id int) (models.App, error) {

	const op = "storage.postgres.GetApp"

	stmt, err := s.db.Prepare("select * from apps where id = $1")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	var app models.App

	err = stmt.QueryRowContext(ctx, id).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil

}

func (s *Storage) GetRoles(ctx context.Context, userId int64) (models.UserRoles, error) {

	const op = "storage.postgres.GetRoles"

	q := "SELECT user_id,role_id,roles.name,roles.app_id FROM user_roles JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = $1"

	rows, err := s.db.QueryContext(ctx, q, userId)

	if err != nil {
		return models.UserRoles{}, fmt.Errorf("%s: %w", op, err)
	}

	var userRoles models.UserRoles

	userRoles.UserId = userId

	for rows.Next() {
		var uid, rid, aid int64
		var roleName string
		if err := rows.Scan(&uid, &rid, &roleName, &aid); err != nil {
			return models.UserRoles{}, fmt.Errorf("%s: scan error: %w", op, err)
		}
		userRoles.Roles = append(userRoles.Roles, models.Role{
			ID:    int(rid),
			Name:  roleName,
			AppId: int(aid),
		})
	}
	if err := rows.Err(); err != nil {
		return models.UserRoles{}, fmt.Errorf("%s: %w", op, err)
	}
	return userRoles, nil
}

func (s *Storage) GetUserRolesByApp(ctx context.Context, userId int64, appId int64) (models.UserRoles, error) {

	const op = "storage.postgres.GetRolesByApp"

	q := "SELECT user_id,role_id,roles.name,roles.app_id FROM user_roles JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = $1 AND roles.app_id = $2"

	rows, err := s.db.QueryContext(ctx, q, userId, appId)

	if err != nil {
		return models.UserRoles{}, fmt.Errorf("%s: %w", op, err)
	}

	var userRoles models.UserRoles

	userRoles.UserId = userId

	for rows.Next() {
		var uid, rid, aid int64
		var roleName string
		if err := rows.Scan(&uid, &rid, &roleName, &aid); err != nil {
			return models.UserRoles{}, fmt.Errorf("%s: scan error: %w", op, err)
		}
		userRoles.Roles = append(userRoles.Roles, models.Role{
			ID:    int(rid),
			Name:  roleName,
			AppId: int(aid),
		})
	}
	if err := rows.Err(); err != nil {
		return models.UserRoles{}, fmt.Errorf("%s: %w", op, err)
	}
	return userRoles, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64, appId int64) (bool, error) {

	const op = "storage.postgres.isAdmin"

	stmt, err := s.db.Prepare("SELECT user_id,role_id,roles.name,roles.app_id FROM user_roles JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = $1 AND roles.app_id = $2 AND roles.name = $3")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRowContext(ctx, userId, appId, "admin").Scan()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotAdmin)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil

}

func (s *Storage) GetRoleByApp(ctx context.Context, appId int64, name string) (models.Role, error) {

	const op = "storage.postgres.GetRoleByApp"

	stmt, err := s.db.Prepare("SELECT * FROM roles WHERE app_id = $1 AND name = $2")
	if err != nil {
		return models.Role{}, fmt.Errorf("%s: %w", op, err)
	}

	var role models.Role

	err = stmt.QueryRowContext(ctx, appId, name).Scan(&role.ID, &role.AppId, &role.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Role{}, fmt.Errorf("%s: %w", op, storage.ErrRoleNotFound)
		}
		return models.Role{}, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil

}

func (s *Storage) SaveRole(
	ctx context.Context,
	userID int64,
	roleID int,
) error {
	const op = "storage.postgres.SaveRole"

	query := `
        INSERT INTO user_roles (user_id, role_id)
        VALUES ($1, $2)
    `
	_, err := s.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		// catch unique-violation: user already has the role
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return storage.ErrUserHaveRole
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
