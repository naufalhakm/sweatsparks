package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sweatsparks/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error
	FindUserByEmail(ctx context.Context, tx *sql.Tx, email string) (*models.User, error)
	FindUserByUsername(ctx context.Context, tx *sql.Tx, username string) (*models.User, error)
	FindUserById(ctx context.Context, tx *sql.Tx, id int) (*models.User, error)
	FindAllUser(ctx context.Context, tx *sql.Tx) ([]*models.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	SQL := "insert into users(username, email, password_hash, created_at, updated_at) values (?, ?, ?, ?, ?)"
	response, err := tx.ExecContext(ctx, SQL, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("Failed to create a user register, transaction rolled back. Reason: " + err.Error())
	}
	userId, err := response.LastInsertId()

	if err != nil {
		return errors.New("Failed to retrieve user_id, transaction rolled back. Reason:" + err.Error())
	}

	user.Id = uint64(userId)

	return nil
}

func (repository *UserRepositoryImpl) FindUserByEmail(ctx context.Context, tx *sql.Tx, email string) (*models.User, error) {
	SQL := "select id, email, username, password_hash, created_at, updated_at from users where email = ?"
	rows, err := tx.QueryContext(ctx, SQL, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("user is not found")
	}
}
func (repository *UserRepositoryImpl) FindUserByUsername(ctx context.Context, tx *sql.Tx, username string) (*models.User, error) {
	SQL := "select id, email, username, password_hash, created_at, updated_at from users where username = ?"
	rows, err := tx.QueryContext(ctx, SQL, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("user is not found")
	}
}

func (repository *UserRepositoryImpl) FindUserById(ctx context.Context, tx *sql.Tx, id int) (*models.User, error) {
	SQL := "select id, email, username, password_hash, created_at, updated_at from users where id = ?"
	rows, err := tx.QueryContext(ctx, SQL, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	} else {
		return nil, errors.New("user is not found")
	}
}

func (repository *UserRepositoryImpl) FindAllUser(ctx context.Context, tx *sql.Tx) ([]*models.User, error) {
	SQL := "select id, email, username, password_hash, created_at, updated_at from users"
	rows, err := tx.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
