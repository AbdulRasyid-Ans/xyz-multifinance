package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
)

type ConsumerLimitRepository interface {
	GetLimitByTenureAndConsumerID(ctx context.Context, tenure int16, consumerID int64) (result ConsumerLimit, err error)
	CreateConsumerLimit(ctx context.Context, consumerLimit ConsumerLimit) (id int64, err error)
	UpdateConsumerLimit(ctx context.Context, consumerLimit ConsumerLimit) (err error)
	DeleteConsumerLimit(ctx context.Context, consumerLimitID int64) (err error)
	GetConsumerLimitByConsumerID(ctx context.Context, consumerID int64) (result []ConsumerLimit, err error)
	GetConsumerLimitByID(ctx context.Context, consumerLimitID int64) (result ConsumerLimit, err error)
}

type consumerLimitRepository struct {
	db *sql.DB
}

func NewConsumerLimitRepository(db *sql.DB) ConsumerLimitRepository {
	return &consumerLimitRepository{db: db}
}

type (
	ConsumerLimit struct {
		ID          int64
		ConsumerID  int64
		Tenure      int16
		LimitAmount float64
	}

	ConsumerLimitScanner struct {
		ID          sql.NullInt64
		ConsumerID  sql.NullInt64
		Tenure      sql.NullInt16
		LimitAmount sql.NullFloat64
	}
)

var (
	ValidConsumerLimitTenure = map[int16]bool{
		1: true,
		2: true,
		3: true,
		6: true,
	}
)

func (r *consumerLimitRepository) GetLimitByTenureAndConsumerID(ctx context.Context, tenure int16, consumerID int64) (result ConsumerLimit, err error) {
	query := `
		SELECT
			consumer_limit_id,
			consumer_id,
			tenure,
			limit_amount
		FROM consumer_limits
		WHERE deleted_at IS NULL
		AND consumer_id = ?
		AND CAST(tenure AS CHAR) = ?
		LIMIT 1
	`

	var consumerLimitScanner ConsumerLimitScanner
	err = r.db.QueryRowContext(ctx, query, consumerID, tenure).Scan(
		&consumerLimitScanner.ID,
		&consumerLimitScanner.ConsumerID,
		&consumerLimitScanner.Tenure,
		&consumerLimitScanner.LimitAmount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		logger.Error(fmt.Sprintf("[consumerLimitRepository][GetLimitByTenureAndConsumerID] while scan query row. Err: %v", err))
		return result, err
	}

	result = ConsumerLimit{
		ID:          consumerLimitScanner.ID.Int64,
		ConsumerID:  consumerLimitScanner.ConsumerID.Int64,
		Tenure:      consumerLimitScanner.Tenure.Int16,
		LimitAmount: consumerLimitScanner.LimitAmount.Float64,
	}

	return result, nil
}

func (r *consumerLimitRepository) CreateConsumerLimit(ctx context.Context, consumerLimit ConsumerLimit) (id int64, err error) {
	query := `
		INSERT INTO consumer_limits (
			consumer_id,
			tenure,
			limit_amount,
			created_at
		) VALUES (?, CAST(? AS CHAR), ?, NOW())
	`

	res, err := r.db.ExecContext(ctx, query,
		consumerLimit.ConsumerID,
		consumerLimit.Tenure,
		consumerLimit.LimitAmount,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("[consumerLimitRepository][CreateConsumerLimit] while exec query. Err: %v", err))
		return 0, err
	}

	return res.LastInsertId()
}

func (r *consumerLimitRepository) UpdateConsumerLimit(ctx context.Context, consumerLimit ConsumerLimit) (err error) {
	query := `
		UPDATE consumer_limits
		SET
			limit_amount = ?,
			updated_at = NOW()
		WHERE deleted_at IS NULL
		AND consumer_limit_id = ?
	`

	_, err = r.db.ExecContext(ctx, query,
		consumerLimit.LimitAmount,
		consumerLimit.ID,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("[consumerLimitRepository][UpdateConsumerLimit] while exec query. Err: %v", err))
		return err
	}

	return nil
}

func (r *consumerLimitRepository) DeleteConsumerLimit(ctx context.Context, consumerLimitID int64) (err error) {
	query := `
		UPDATE consumer_limits
		SET
			deleted_at = NOW()
		WHERE consumer_limit_id = ?
	`

	_, err = r.db.ExecContext(ctx, query, consumerLimitID)
	if err != nil {
		logger.Error(fmt.Sprintf("[consumerLimitRepository][DeleteConsumerLimit] while exec query. Err: %v", err))
		return err
	}

	return nil
}

func (r *consumerLimitRepository) GetConsumerLimitByConsumerID(ctx context.Context, consumerID int64) (result []ConsumerLimit, err error) {
	query := `
		SELECT
			consumer_limit_id,
			consumer_id,
			tenure,
			limit_amount
		FROM consumer_limits
		WHERE deleted_at IS NULL
		AND consumer_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, consumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[consumerLimitRepository][GetConsumerLimitByConsumerID] while query. Err: %v", err))
		return nil, err
	}
	defer rows.Close()

	var consumerLimits []ConsumerLimit
	for rows.Next() {
		var consumerLimitScanner ConsumerLimitScanner
		err := rows.Scan(
			&consumerLimitScanner.ID,
			&consumerLimitScanner.ConsumerID,
			&consumerLimitScanner.Tenure,
			&consumerLimitScanner.LimitAmount,
		)
		if err != nil {
			logger.Error(fmt.Sprintf("[consumerLimitRepository][GetConsumerLimitByConsumerID] while scan query row. Err: %v", err))
			return nil, err
		}

		consumerLimit := ConsumerLimit{
			ID:          consumerLimitScanner.ID.Int64,
			ConsumerID:  consumerLimitScanner.ConsumerID.Int64,
			Tenure:      consumerLimitScanner.Tenure.Int16,
			LimitAmount: consumerLimitScanner.LimitAmount.Float64,
		}

		consumerLimits = append(consumerLimits, consumerLimit)
	}

	return consumerLimits, nil
}

func (r *consumerLimitRepository) GetConsumerLimitByID(ctx context.Context, consumerLimitID int64) (result ConsumerLimit, err error) {
	query := `
		SELECT
			consumer_limit_id,
			consumer_id,
			tenure,
			limit_amount
		FROM consumer_limits
		WHERE deleted_at IS NULL
		AND consumer_limit_id = ?
		LIMIT 1
	`

	var consumerLimitScanner ConsumerLimitScanner
	err = r.db.QueryRowContext(ctx, query, consumerLimitID).Scan(
		&consumerLimitScanner.ID,
		&consumerLimitScanner.ConsumerID,
		&consumerLimitScanner.Tenure,
		&consumerLimitScanner.LimitAmount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		logger.Error(fmt.Sprintf("[consumerLimitRepository][GetConsumerLimitByID] while scan query row. Err: %v", err))
		return result, err
	}

	result = ConsumerLimit{
		ID:          consumerLimitScanner.ID.Int64,
		ConsumerID:  consumerLimitScanner.ConsumerID.Int64,
		Tenure:      consumerLimitScanner.Tenure.Int16,
		LimitAmount: consumerLimitScanner.LimitAmount.Float64,
	}

	return result, nil
}
