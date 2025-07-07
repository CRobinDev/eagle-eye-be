package repository

import (
	"context"

	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	conn *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) interfaces.IUserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (ur *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
	query, args, err := squirrel.
		Select("*").
		From("users").
		Where("id = ? AND deleted_at IS NULL", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.User{}, err
	}

	var user entities.User
	err = ur.conn.GetContext(ctx, &user, query, args...)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	queryBuilder := squirrel.
		Update("users")

	if user.IsCustomer {
		queryBuilder = queryBuilder.Set("is_customer", true)
	}

	if user.Username != "" {
		queryBuilder = queryBuilder.Set("username", user.Username)
	}

	query, values, err := queryBuilder.
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = ur.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) UpdatePassword(ctx context.Context, user *entities.User) error {
	query, values, err := squirrel.
		Update("users").
		Set("password_hash", user.Password).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = ur.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.
		Update("users").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where("id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = ur.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
