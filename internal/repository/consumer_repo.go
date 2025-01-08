package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ConsumerRepository interface {
	CreateConsumer(ctx context.Context, consumer Consumer) (id int64, err error)
	GetConsumerByID(ctx context.Context, id int64) (data Consumer, err error)
	FetchConsumer(ctx context.Context, req FetchConsumerRequest) (data []Consumer, err error)
	UpdateConsumer(ctx context.Context, consumer Consumer) (err error)
	DeleteConsumer(ctx context.Context, id int64) (err error)
}

type consumerRepo struct {
	db *sql.DB
}

func NewConsumerRepository(db *sql.DB) ConsumerRepository {
	return &consumerRepo{db: db}
}

type (
	FetchConsumerRequest struct {
		Limit  int
		Offset int
	}

	Consumer struct {
		ID           int64
		FullName     string
		LegalName    string
		PlaceOfBirth string
		DOB          string
		Salary       float64
		NIK          string
		KTPImageURL  string
		SelfieURL    string
		CreatedAt    time.Time
	}
)

type ConsumerScanner struct {
	ID           sql.NullInt64
	FullName     sql.NullString
	LegalName    sql.NullString
	PlaceOfBirth sql.NullString
	DOB          sql.NullString
	Salary       sql.NullFloat64
	NIK          sql.NullString
	KTPImageURL  sql.NullString
	SelfieURL    sql.NullString
	CreatedAt    sql.NullTime
}

func (r *consumerRepo) CreateConsumer(ctx context.Context, consumer Consumer) (id int64, err error) {
	query := `
		INSERT INTO consumers (
			full_name,
			legal_name,
			place_of_birth,
			dob,
			salary,
			nik,
			ktp_image_url,
			selfie_image_url,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		consumer.FullName,
		consumer.LegalName,
		consumer.PlaceOfBirth,
		consumer.DOB,
		consumer.Salary,
		consumer.NIK,
		consumer.KTPImageURL,
		consumer.SelfieURL,
	)
	if err != nil {
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *consumerRepo) UpdateConsumer(ctx context.Context, consumer Consumer) (err error) {
	query := `
		UPDATE consumers
		SET
			full_name = ?,
			legal_name = ?,
			place_of_birth = ?,
			dob = ?,
			salary = ?,
			nik = ?,
			ktp_image_url = ?,
			selfie_image_url = ?,
			modified_at = NOW()
		WHERE deleted_at is null
		AND consumer_id = ?
	`

	_, err = r.db.ExecContext(ctx, query,
		consumer.FullName,
		consumer.LegalName,
		consumer.PlaceOfBirth,
		consumer.DOB,
		consumer.Salary,
		consumer.NIK,
		consumer.KTPImageURL,
		consumer.SelfieURL,
		consumer.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *consumerRepo) DeleteConsumer(ctx context.Context, id int64) (err error) {
	query := `
		UPDATE consumers
		SET
			deleted_at = NOW()
		WHERE consumer_id = ?
	`

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *consumerRepo) GetConsumerByID(ctx context.Context, id int64) (data Consumer, err error) {
	query := `
		SELECT
			consumer_id,
			full_name,
			legal_name,
			place_of_birth,
			dob,
			salary,
			nik,
			ktp_image_url,
			selfie_image_url,
			created_at
		FROM consumers
		WHERE deleted_at is null
		AND consumer_id = ?
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var consumerScanner ConsumerScanner
	if err := row.Scan(
		&consumerScanner.ID,
		&consumerScanner.FullName,
		&consumerScanner.LegalName,
		&consumerScanner.PlaceOfBirth,
		&consumerScanner.DOB,
		&consumerScanner.Salary,
		&consumerScanner.NIK,
		&consumerScanner.KTPImageURL,
		&consumerScanner.SelfieURL,
		&consumerScanner.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data, nil
		}
		return data, err
	}

	data = Consumer{
		ID:           consumerScanner.ID.Int64,
		FullName:     consumerScanner.FullName.String,
		LegalName:    consumerScanner.LegalName.String,
		PlaceOfBirth: consumerScanner.PlaceOfBirth.String,
		DOB:          consumerScanner.DOB.String,
		Salary:       consumerScanner.Salary.Float64,
		NIK:          consumerScanner.NIK.String,
		KTPImageURL:  consumerScanner.KTPImageURL.String,
		SelfieURL:    consumerScanner.SelfieURL.String,
		CreatedAt:    consumerScanner.CreatedAt.Time,
	}

	return data, nil
}

func (r *consumerRepo) FetchConsumer(ctx context.Context, req FetchConsumerRequest) (data []Consumer, err error) {
	query := `
		SELECT
			consumer_id,
			full_name,
			legal_name,
			place_of_birth,
			dob,
			salary,
			nik,
			ktp_image_url,
			selfie_image_url,
			created_at
		FROM consumers
		WHERE deleted_at is null
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var consumerScanner ConsumerScanner
		if err := rows.Scan(
			&consumerScanner.ID,
			&consumerScanner.FullName,
			&consumerScanner.LegalName,
			&consumerScanner.PlaceOfBirth,
			&consumerScanner.DOB,
			&consumerScanner.Salary,
			&consumerScanner.NIK,
			&consumerScanner.KTPImageURL,
			&consumerScanner.SelfieURL,
			&consumerScanner.CreatedAt,
		); err != nil {
			return nil, err
		}

		data = append(data, Consumer{
			ID:           consumerScanner.ID.Int64,
			FullName:     consumerScanner.FullName.String,
			LegalName:    consumerScanner.LegalName.String,
			PlaceOfBirth: consumerScanner.PlaceOfBirth.String,
			DOB:          consumerScanner.DOB.String,
			Salary:       consumerScanner.Salary.Float64,
			NIK:          consumerScanner.NIK.String,
			KTPImageURL:  consumerScanner.KTPImageURL.String,
			SelfieURL:    consumerScanner.SelfieURL.String,
			CreatedAt:    consumerScanner.CreatedAt.Time,
		})
	}

	return data, nil
}
