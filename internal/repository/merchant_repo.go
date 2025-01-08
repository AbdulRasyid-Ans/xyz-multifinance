package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type MerchantRepository interface {
	CreateMerchant(ctx context.Context, merchant Merchant) (id int64, err error)
	GetMerchantByID(ctx context.Context, id int64) (data Merchant, err error)
	FetchMerchant(ctx context.Context, req FetchMerchantRequest) (data []Merchant, err error)
	UpdateMerchant(ctx context.Context, merchant Merchant) (err error)
	DeleteMerchant(ctx context.Context, id int64) (err error)
}

type merchantRepo struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) MerchantRepository {
	return &merchantRepo{db: db}
}

type (
	FetchMerchantRequest struct {
		Limit  int
		Offset int
	}

	Merchant struct {
		ID           int64
		MerchantName string
		MerchantType string
		CreatedAt    time.Time
	}
)

type MerchantScanner struct {
	ID           sql.NullInt64
	MerchantName sql.NullString
	MerchantType sql.NullString
	CreatedAt    sql.NullTime
}

func (r *merchantRepo) CreateMerchant(ctx context.Context, merchant Merchant) (id int64, err error) {
	query := `
		INSERT INTO merchants (
			merchant_name, 
			merchant_type, 
			created_at
		) VALUES (?, ?, NOW())
	`

	res, err := r.db.ExecContext(ctx, query,
		merchant.MerchantName,
		merchant.MerchantType,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *merchantRepo) GetMerchantByID(ctx context.Context, id int64) (data Merchant, err error) {
	query := `
		SELECT 
			merchant_id, 
			merchant_name, 
			merchant_type, 
			created_at 
		FROM merchants 
		WHERE deleted_at IS NULL 
		AND merchant_id = ?
		LIMIT 1
	`

	var merchantScanner MerchantScanner
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&merchantScanner.ID,
		&merchantScanner.MerchantName,
		&merchantScanner.MerchantType,
		&merchantScanner.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data, nil
		}

		return data, err
	}

	data = Merchant{
		ID:           merchantScanner.ID.Int64,
		MerchantName: merchantScanner.MerchantName.String,
		MerchantType: merchantScanner.MerchantType.String,
		CreatedAt:    merchantScanner.CreatedAt.Time,
	}

	return data, nil
}

func (r *merchantRepo) FetchMerchant(ctx context.Context, req FetchMerchantRequest) (data []Merchant, err error) {
	query := `
		SELECT 
			merchant_id, 
			merchant_name, 
			merchant_type, 
			created_at 
		FROM merchants 
		WHERE deleted_at IS NULL 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var merchantScanner MerchantScanner
		err = rows.Scan(
			&merchantScanner.ID,
			&merchantScanner.MerchantName,
			&merchantScanner.MerchantType,
			&merchantScanner.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		merchant := Merchant{
			ID:           merchantScanner.ID.Int64,
			MerchantName: merchantScanner.MerchantName.String,
			MerchantType: merchantScanner.MerchantType.String,
			CreatedAt:    merchantScanner.CreatedAt.Time,
		}

		data = append(data, merchant)
	}

	return data, nil
}

func (r *merchantRepo) UpdateMerchant(ctx context.Context, merchant Merchant) (err error) {
	query := `
		UPDATE merchants 
		SET 
			merchant_name = ?, 
			merchant_type = ?
		WHERE deleted_at is null 
		AND merchant_id = ?
	`

	_, err = r.db.ExecContext(ctx, query,
		merchant.MerchantName,
		merchant.MerchantType,
		merchant.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *merchantRepo) DeleteMerchant(ctx context.Context, id int64) (err error) {
	query := `
		UPDATE merchants 
		SET 
			deleted_at = NOW() 
		WHERE merchant_id = ?
	`

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
