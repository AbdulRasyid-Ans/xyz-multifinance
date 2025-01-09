package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan Loan) (int64, error)
	UpdateLoan(ctx context.Context, req UpdateLoanRequest, tx *sql.Tx) error
	GetLoanByID(ctx context.Context, loanID int64) (Loan, error)
	DeleteLoan(ctx context.Context, loanID int64) error
	GetLoanByConsumerID(ctx context.Context, consumerID int64) ([]Loan, error)
}

type loanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepository{db: db}
}

type (
	UpdateLoanRequest struct {
		ID                 int64
		PaidLoanAmount     float64
		PaidInterestAmount float64
		LoanStatus         string
		Installment        int32
	}

	Loan struct {
		ID                 int64
		ConsumerLimitID    int64
		ConsumerID         int64
		MerchantID         int64
		LoanAmount         float64
		PaidLoanAmount     float64
		ContractNumber     string
		InterestRate       float64
		InterestAmount     float64
		PaidInterestAmount float64
		LoanStatus         string
		DueDate            time.Time
		Installment        int32
		AssetName          string
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}

	LoanScanner struct {
		ID                 sql.NullInt64
		ConsumerLimitID    sql.NullInt64
		ConsumerID         sql.NullInt64
		MerchantID         sql.NullInt64
		LoanAmount         sql.NullFloat64
		PaidLoanAmount     sql.NullFloat64
		ContractNumber     sql.NullString
		InterestRate       sql.NullFloat64
		InterestAmount     sql.NullFloat64
		PaidInterestAmount sql.NullFloat64
		LoanStatus         sql.NullString
		DueDate            sql.NullTime
		Installment        sql.NullInt32
		AssetName          sql.NullString
		CreatedAt          sql.NullTime
		UpdatedAt          sql.NullTime
	}
)

var (
	ValidLoanStatus = map[string]bool{
		"on_going": true,
		"finish":   true,
		"late":     true,
	}
)

func (r *loanRepository) CreateLoan(ctx context.Context, loan Loan) (id int64, err error) {
	query := `
		INSERT INTO loans (
			consumer_limit_id,
			consumer_id,
			merchant_id,
			loan_amount,
			contract_number,
			interest_rate,
			interest_amount,
			due_date,
			asset_name,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		loan.ConsumerLimitID,
		loan.ConsumerID,
		loan.MerchantID,
		loan.LoanAmount,
		loan.ContractNumber,
		loan.InterestRate,
		loan.InterestAmount,
		loan.DueDate,
		loan.AssetName,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("[loanRepository][CreateLoan] while exec query. Err: %v", err))
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		logger.Error(fmt.Sprintf("[loanRepository][CreateLoan] while get last insert id. Err: %v", err))
		return id, err
	}

	return id, nil
}

func (r *loanRepository) UpdateLoan(ctx context.Context, req UpdateLoanRequest, tx *sql.Tx) (err error) {
	query := `
		UPDATE loans
		SET
			paid_loan_amount = ?,
			paid_interest_amount = ?,
			loan_status = ?,
			installment = ?,
			updated_at = NOW()
		WHERE deleted_at IS NULL
		AND loan_id = ?
	`

	if tx != nil {
		_, err = tx.ExecContext(ctx, query,
			req.PaidLoanAmount,
			req.PaidInterestAmount,
			req.LoanStatus,
			req.Installment,
			req.ID,
		)
	} else {
		_, err = r.db.ExecContext(ctx, query,
			req.PaidLoanAmount,
			req.PaidInterestAmount,
			req.LoanStatus,
			req.Installment,
			req.ID,
		)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("[loanRepository][UpdateLoan] while exec query. Err: %v", err))
		return err
	}

	return nil
}

func (r *loanRepository) GetLoanByID(ctx context.Context, loanID int64) (result Loan, err error) {
	query := `
		SELECT
			loan_id,
			consumer_limit_id,
			consumer_id,
			merchant_id,
			loan_amount,
			paid_loan_amount,
			contract_number,
			interest_rate,
			interest_amount,
			paid_interest_amount,
			loan_status,
			due_date,
			installment,
			asset_name,
			created_at,
			updated_at
		FROM loans
		WHERE deleted_at IS NULL
		AND loan_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, loanID)

	var loanScanner LoanScanner
	err = row.Scan(
		&loanScanner.ID,
		&loanScanner.ConsumerLimitID,
		&loanScanner.ConsumerID,
		&loanScanner.MerchantID,
		&loanScanner.LoanAmount,
		&loanScanner.PaidLoanAmount,
		&loanScanner.ContractNumber,
		&loanScanner.InterestRate,
		&loanScanner.InterestAmount,
		&loanScanner.PaidInterestAmount,
		&loanScanner.LoanStatus,
		&loanScanner.DueDate,
		&loanScanner.Installment,
		&loanScanner.AssetName,
		&loanScanner.CreatedAt,
		&loanScanner.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		logger.Error(fmt.Sprintf("[loanRepository][GetLoanByID] while scan query row. Err: %v", err))
		return result, err
	}

	result = Loan{
		ID:                 loanScanner.ID.Int64,
		ConsumerLimitID:    loanScanner.ConsumerLimitID.Int64,
		ConsumerID:         loanScanner.ConsumerID.Int64,
		MerchantID:         loanScanner.MerchantID.Int64,
		LoanAmount:         loanScanner.LoanAmount.Float64,
		PaidLoanAmount:     loanScanner.PaidLoanAmount.Float64,
		ContractNumber:     loanScanner.ContractNumber.String,
		InterestRate:       loanScanner.InterestRate.Float64,
		InterestAmount:     loanScanner.InterestAmount.Float64,
		PaidInterestAmount: loanScanner.PaidInterestAmount.Float64,
		LoanStatus:         loanScanner.LoanStatus.String,
		DueDate:            loanScanner.DueDate.Time,
		Installment:        loanScanner.Installment.Int32,
		AssetName:          loanScanner.AssetName.String,
		CreatedAt:          loanScanner.CreatedAt.Time,
		UpdatedAt:          loanScanner.UpdatedAt.Time,
	}

	return result, nil
}

func (r *loanRepository) DeleteLoan(ctx context.Context, loanID int64) (err error) {
	query := `
		UPDATE loans
		SET
			deleted_at = NOW()
		WHERE loan_id = ?
	`

	_, err = r.db.ExecContext(ctx, query, loanID)
	if err != nil {
		logger.Error(fmt.Sprintf("[loanRepository][DeleteLoan] while exec query. Err: %v", err))
		return err
	}

	return nil
}

func (r *loanRepository) GetLoanByConsumerID(ctx context.Context, consumerID int64) (result []Loan, err error) {
	query := `
		SELECT
			loan_id,
			consumer_limit_id,
			consumer_id,
			merchant_id,
			loan_amount,
			paid_loan_amount,
			contract_number,
			interest_rate,
			interest_amount,
			paid_interest_amount,
			loan_status,
			due_date,
			installment,
			asset_name,
			created_at,
			updated_at
		FROM loans
		WHERE deleted_at IS NULL
		AND consumer_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, consumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[loanRepository][GetLoanByConsumerID] while query. Err: %v", err))
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var loanScanner LoanScanner
		err = rows.Scan(
			&loanScanner.ID,
			&loanScanner.ConsumerLimitID,
			&loanScanner.ConsumerID,
			&loanScanner.MerchantID,
			&loanScanner.LoanAmount,
			&loanScanner.PaidLoanAmount,
			&loanScanner.ContractNumber,
			&loanScanner.InterestRate,
			&loanScanner.InterestAmount,
			&loanScanner.PaidInterestAmount,
			&loanScanner.LoanStatus,
			&loanScanner.DueDate,
			&loanScanner.Installment,
			&loanScanner.AssetName,
			&loanScanner.CreatedAt,
			&loanScanner.UpdatedAt,
		)
		if err != nil {
			logger.Error(fmt.Sprintf("[loanRepository][GetLoanByConsumerID] while scan query row. Err: %v", err))
			return result, err
		}

		loan := Loan{
			ID:                 loanScanner.ID.Int64,
			ConsumerLimitID:    loanScanner.ConsumerLimitID.Int64,
			ConsumerID:         loanScanner.ConsumerID.Int64,
			MerchantID:         loanScanner.MerchantID.Int64,
			LoanAmount:         loanScanner.LoanAmount.Float64,
			PaidLoanAmount:     loanScanner.PaidLoanAmount.Float64,
			ContractNumber:     loanScanner.ContractNumber.String,
			InterestRate:       loanScanner.InterestRate.Float64,
			InterestAmount:     loanScanner.InterestAmount.Float64,
			PaidInterestAmount: loanScanner.PaidInterestAmount.Float64,
			LoanStatus:         loanScanner.LoanStatus.String,
			DueDate:            loanScanner.DueDate.Time,
			Installment:        loanScanner.Installment.Int32,
			AssetName:          loanScanner.AssetName.String,
			CreatedAt:          loanScanner.CreatedAt.Time,
			UpdatedAt:          loanScanner.UpdatedAt.Time,
		}

		result = append(result, loan)
	}

	return result, nil
}
