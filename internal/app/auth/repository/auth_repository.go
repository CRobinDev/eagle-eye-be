package repository

import (
	"context"

	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type authRepository struct {
	conn *sqlx.DB
}

func NewAuthRepository(conn *sqlx.DB) interfaces.IAuthRepository {
	return &authRepository{
		conn: conn,
	}
}

func (ar *authRepository) CreateUser(ctx context.Context, user *entities.User) error {
	tx, err := ar.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	query, values, err := squirrel.
		Insert("users").
		Columns("id", "username", "email", "password_hash").
		Values(user.ID, user.Username, user.Email, user.Password).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				tx.Rollback()
				return err

			}
		}

		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	query, args, err := squirrel.
		Select("id", "username", "role", "is_customer", "password_hash", "created_at").
		From("users").
		Where(squirrel.Eq{
			"email":      email,
			"deleted_at": nil,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.User{}, err
	}

	var user entities.User
	err = ar.conn.GetContext(ctx, &user, query, args...)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (ar *authRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (entities.User, error) {
	query, args, err := squirrel.
		Select("id", "username", "email", "is_verified", "is_customer").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.User{}, err
	}

	var user entities.User

	err = ar.conn.GetContext(ctx, &user, query, args...)
	if err != nil {
		return entities.User{}, err
	}

	return user, err
}
