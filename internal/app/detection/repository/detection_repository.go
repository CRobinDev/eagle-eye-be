package repository

import (
	"context"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type detectionRepository struct {
	conn *sqlx.DB
}

func NewDetectionRepository(conn *sqlx.DB) interfaces.IDetectionRepository {
	return &detectionRepository{
		conn: conn,
	}
}

func (dr *detectionRepository) CreateDetection(ctx context.Context, detection *entities.Detection) error {
	query, values, err := squirrel.
		Insert("detections").
		Columns("ip_address", "email", "customer_id", "path", "method", "status_code").
		Values(detection.IPAddress, detection.Email, detection.CustomerID, detection.Path, detection.Method, detection.StatusCode).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	query = dr.conn.Rebind(query)
	_, err = dr.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil

}

func (dr *detectionRepository) GetDetection(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error) {
	queryBuilder := squirrel.Select("id", "ip_address", "path", "status_code", "method", "email", "created_at").
		From("detections").
		Where(squirrel.Expr("deleted_at IS NULL"))

	if filter.CustomerID != uuid.Nil {
		queryBuilder = queryBuilder.
			Where(squirrel.Eq{
				"customer_id": filter.CustomerID,
			})
	}

	query, args, err := queryBuilder.
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var detections []entities.Detection
	err = dr.conn.SelectContext(ctx, &detections, query, args...)
	if err != nil {
		return nil, err
	}

	return detections, nil
}

func (dr *detectionRepository) DeleteDetection(ctx context.Context, id uint16) error {
	query, values, err := squirrel.
		Update("detections").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{
			"id": id,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	query = dr.conn.Rebind(query)
	_, err = dr.conn.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (dr *detectionRepository) GetDeepFakeDetected(ctx context.Context, filter dto.GetDetectionFilter) ([]entities.Detection, error) {
	queryBuilder := squirrel.Select("id", "ip_address", "path", "method", "status_code", "email", "created_at").
		From("detections").
		Where(squirrel.Eq{"is_deepfake": true}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	if filter.CustomerID != uuid.Nil {
		queryBuilder = queryBuilder.
			Where(squirrel.Eq{
				"customer_id": filter.CustomerID,
			})
	}

	query, args, err := queryBuilder.
		OrderBy("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var detections []entities.Detection
	err = dr.conn.SelectContext(ctx, &detections, query, args...)
	if err != nil {
		return nil, err
	}

	return detections, nil
}
