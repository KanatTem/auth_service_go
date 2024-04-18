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

	stmt, err := s.db.Prepare("INSERT into users(email, password_hash) VALUES ($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			// 23505 = unique_violation
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
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
