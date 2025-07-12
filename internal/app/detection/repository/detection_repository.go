package repository

import (
	"context"
	"time"

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
	queryBuilder := squirrel.Select("id", "ip_address", "path", "status_code", "type", "method", "email", "is_deepfake", "created_at").
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

func (dr *detectionRepository) GetDetectionByID(ctx context.Context, id uint16) (entities.Detection, error) {
	query, args, err := squirrel.
		Select("*").
		From("detections").
		Where(squirrel.Eq{
			"id": id,
		}).
		OrderBy("created_at DESC").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return entities.Detection{}, err
	}

	var detection entities.Detection
	err = dr.conn.GetContext(ctx, &detection, query, args...)
	if err != nil {
		return entities.Detection{}, err
	}

	return detection, err
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

func (dr *detectionRepository) UndeleteDetection(ctx context.Context, id uint16) error {
	query, values, err := squirrel.
		Update("detections").
		Set("deleted_at", nil).
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
	queryBuilder := squirrel.Select("id", "type", "ip_address", "path", "method", "status_code", "email", "is_deepfake", "created_at").
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

func (dr *detectionRepository) GetHourlyUsage(ctx context.Context, customerID uuid.UUID, date time.Time) ([]dto.CustomerUsage, error) {

	query, args, err := squirrel.
		Select("TO_CHAR(DATE_TRUNC('hour', created_at), 'HH24:00') AS time", "COUNT(*) AS usage").
		From("detections").
		Where(squirrel.Eq{"customer_id": customerID}).
		Where(squirrel.Expr("DATE(created_at) = ?", date.Format("2006-01-02"))).
		GroupBy("time").
		OrderBy("time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var usages []dto.CustomerUsage
	err = dr.conn.SelectContext(ctx, &usages, query, args...)
	if err != nil {
		return nil, err
	}

	return usages, nil
}

func (dr *detectionRepository) GetDailyUsage(ctx context.Context, customerID uuid.UUID, days int) ([]dto.CustomerUsage, error) {

	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	query, args, err := squirrel.
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') AS time", "COUNT(*) AS usage").
		From("detections").
		Where(squirrel.Eq{"customer_id": customerID}).
		Where(squirrel.Expr("created_at >= ? AND created_at <= NOW()", startDate)).
		GroupBy("time").
		OrderBy("time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var usages []dto.CustomerUsage
	err = dr.conn.SelectContext(ctx, &usages, query, args...)
	if err != nil {
		return nil, err
	}

	return usages, nil
}

func (dr *detectionRepository) GetWeeklyUsage(ctx context.Context, customerID uuid.UUID, weeks int) ([]dto.CustomerUsage, error) {

	startDate := time.Now().AddDate(0, 0, -weeks*7).Format("2006-01-02")

	query, args, err := squirrel.
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') AS time", "COUNT(*) AS usage").
		From("detections").
		Where(squirrel.Eq{"customer_id": customerID}).
		Where(squirrel.Expr("created_at >= ? AND created_at <= NOW()", startDate)).
		GroupBy("time").
		OrderBy("time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var usages []dto.CustomerUsage
	err = dr.conn.SelectContext(ctx, &usages, query, args...)
	if err != nil {
		return nil, err
	}

	return usages, nil
}

func (dr *detectionRepository) GetMonthlyUsage(ctx context.Context, customerID uuid.UUID, months int) ([]dto.CustomerUsage, error) {

	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01-02")

	query, args, err := squirrel.
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') AS time", "COUNT(*) AS usage").
		From("detections").
		Where(squirrel.Eq{"customer_id": customerID}).
		Where(squirrel.Expr("created_at >= ? AND created_at <= NOW()", startDate)).
		GroupBy("time").
		OrderBy("time").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var usages []dto.CustomerUsage
	err = dr.conn.SelectContext(ctx, &usages, query, args...)
	if err != nil {
		return nil, err
	}

	return usages, nil
}
