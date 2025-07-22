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

type customerRepository struct {
	conn *sqlx.DB
}

func NewCustomerRepository(conn *sqlx.DB) interfaces.ICustomerRepository {
	return &customerRepository{
		conn: conn,
	}
}

func (cr *customerRepository) CreateCustomer(ctx context.Context, customer *entities.Customer) error {
	tx, err := cr.conn.BeginTx(ctx, nil)
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
		Insert("customers").
		Columns("id", "order_id", "customer_tier", "hashed_key", "prefix", "monthly_limit", "expires_at").
		Values(customer.ID, customer.OrderID, customer.CustomerTier, customer.ApiKey, customer.Prefix, customer.MonthlyLimit, customer.ExpiresAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
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

func (cr *customerRepository) UpdateAPIKey(ctx context.Context, customer *entities.Customer) (entities.Customer, error) {
	query, values, err := squirrel.
		Update("customers").
		Set("hashed_key", customer.ApiKey).
		Set("prefix", customer.Prefix).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": customer.ID}).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING expires_at").
		ToSql()

	if err != nil {
		return entities.Customer{}, err
	}

	var suffixCustomer entities.Customer
	err = cr.conn.GetContext(ctx, &suffixCustomer, query, values...)
	if err != nil {
		return entities.Customer{}, err
	}

	return suffixCustomer, nil
}

func (cr *customerRepository) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (entities.Customer, error) {
	query, args, err := squirrel.
		Select("*").
		From("customers").
		Where(
			squirrel.And{
				squirrel.Eq{
					"id":         customerID,
					"revoked_at": nil,
				},
				squirrel.Expr("expires_at > NOW()"),
			}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.Customer{}, err
	}

	var customer entities.Customer
	err = cr.conn.GetContext(ctx, &customer, query, args...)
	if err != nil {
		return entities.Customer{}, err
	}

	return customer, nil
}

func (cr *customerRepository) UpdateUsage(ctx context.Context, prefix string) (entities.Customer, error) {
	query, values, err := squirrel.
		Update("customers").
		Set("current_usage", squirrel.Expr("current_usage + 1")).
		Set("last_used", squirrel.Expr("NOW()")).
		Where(
			squirrel.And{
				squirrel.Eq{
					"prefix":     prefix,
					"revoked_at": nil,
				},
				squirrel.Expr("expires_at > NOW()"),
			},
		).
		Suffix("RETURNING id, customer_tier, hashed_key, prefix, current_usage, monthly_limit, last_used").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.Customer{}, err
	}

	var customer entities.Customer
	err = cr.conn.GetContext(ctx, &customer, query, values...)
	if err != nil {
		return entities.Customer{}, err
	}

	return customer, nil
}

func (cr *customerRepository) GetCustomerTier(ctx context.Context, customerID uuid.UUID) (entities.CustomerTier, error) {
	query, args, err := squirrel.
		Select("customer_tier").
		From("customers").
		Where(squirrel.Eq{
			"id": customerID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, err
	}

	var tier entities.CustomerTier
	err = cr.conn.GetContext(ctx, &tier, query, args...)
	if err != nil {
		return 0, err
	}

	return tier, nil
}
