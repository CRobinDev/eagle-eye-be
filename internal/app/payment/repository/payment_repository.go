package repository

import (
	"context"
	"time"

	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type paymentRepository struct {
	conn *sqlx.DB
}

func NewPaymentRepository(conn *sqlx.DB) interfaces.IPaymentRepository {
	return &paymentRepository{
		conn: conn,
	}
}

func (pr *paymentRepository) CreatePayment(ctx context.Context, payment *entities.Payment) error {
	tx, err := pr.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	queryBuilder := squirrel.
		Insert("payments")

	if payment.Tier != "" {
		queryBuilder = queryBuilder.
			Columns("user_id", "order_id", "amount", "tier").
			Values(payment.UserID, payment.OrderID, payment.Amount, payment.Tier)
	} else {
		queryBuilder = queryBuilder.
			Columns("user_id", "order_id", "amount").
			Values(payment.UserID, payment.OrderID, payment.Amount)
	}

	query, values, err := queryBuilder.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (pr *paymentRepository) UpdatePaymentStatus(ctx context.Context, payment *entities.Payment) error {
	tx, err := pr.conn.BeginTx(ctx, nil)
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
		Update("payments").
		Set("status", payment.Status).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"order_id": payment.OrderID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		tx.Rollback()
		return err
	}

	query = pr.conn.Rebind(query)
	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (pr *paymentRepository) GetStatusByOrderID(ctx context.Context, orderID string) (entities.Payment, error) {
	query, args, err := squirrel.
		Select("*").
		From("payments").
		Where(squirrel.Eq{
			"order_id": orderID,
			"status":   "pending",
		}).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.Payment{}, err
	}

	query = pr.conn.Rebind(query)
	var payment entities.Payment
	if err = pr.conn.GetContext(ctx, &payment, query, args); err != nil {
		return entities.Payment{}, err
	}

	return payment, nil
}

func (pr *paymentRepository) GetStatusByUserID(ctx context.Context, userID uuid.UUID) (entities.Payment, error) {
	query, args, err := squirrel.
		Select("status", "order_id").
		From("payments").
		Where(squirrel.Eq{
			"user_id": userID,
		}).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.Payment{}, err
	}

	query = pr.conn.Rebind(query)

	var payment entities.Payment
	if err = pr.conn.GetContext(ctx, &payment, query, args...); err != nil {
		return entities.Payment{}, err
	}

	return payment, nil
}
