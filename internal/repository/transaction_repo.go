package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/logger"
)

type TransactionRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(ctx context.Context, tx *sql.Tx) error
	RollbackTx(ctx context.Context, tx *sql.Tx) error
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction Transaction) (id int64, err error)
	GetTransactionByID(ctx context.Context, id int64) (result Transaction, err error)
	GetTransactionsByConsumerID(ctx context.Context, consumerID int64) (results []Transaction, err error)
	GetTransactionsByLoanID(ctx context.Context, loanID int64) (results []Transaction, err error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][BeginTx] while begin sql transaction. Err: %v", err))
		return nil, err
	}
	return tx, nil
}

func (r *transactionRepository) CommitTx(ctx context.Context, tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][CommitTx] while commit sql transaction. Err: %v", err))
		return err
	}
	return nil
}

func (r *transactionRepository) RollbackTx(ctx context.Context, tx *sql.Tx) error {
	err := tx.Rollback()
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][RollbackTx] while rollback sql transaction. Err: %v", err))
		return err
	}
	return nil
}

type (
	Transaction struct {
		ID          int64
		ConsumerID  int64
		LoanID      int64
		Amount      float64
		Description string
		CreatedAt   time.Time
	}

	TransactionScanner struct {
		ID          sql.NullInt64
		ConsumerID  sql.NullInt64
		LoanID      sql.NullInt64
		Amount      sql.NullFloat64
		Description sql.NullString
		CreatedAt   sql.NullTime
	}
)

func (r *transactionRepository) CreateTransaction(ctx context.Context, tx *sql.Tx, transaction Transaction) (id int64, err error) {
	query := `
		INSERT INTO transactions (
			consumer_id,
			loan_id,
			amount,
			description,
			created_at
		)
		VALUES (?, ?, ?, ?, now())
	`

	result, err := tx.ExecContext(ctx, query,
		transaction.ConsumerID,
		transaction.LoanID,
		transaction.Amount,
		transaction.Description,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][CreateTransaction] while exec query. Err: %v", err))
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][CreateTransaction] while get last insert id. Err: %v", err))
		return id, err
	}

	return id, nil
}

func (r *transactionRepository) GetTransactionByID(ctx context.Context, id int64) (result Transaction, err error) {
	query := `
		SELECT
			transaction_id,
			consumer_id,
			loan_id,
			amount,
			description,
			created_at
		FROM transactions
		WHERE deleted_at IS NULL
		AND transaction_id = ?
	`

	var transactionScanner TransactionScanner
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&transactionScanner.ID,
		&transactionScanner.ConsumerID,
		&transactionScanner.LoanID,
		&transactionScanner.Amount,
		&transactionScanner.Description,
		&transactionScanner.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		logger.Error(fmt.Sprintf("[transactionRepository][GetTransactionByID] while scan query row. Err: %v", err))
		return result, err
	}

	result = Transaction{
		ID:          transactionScanner.ID.Int64,
		ConsumerID:  transactionScanner.ConsumerID.Int64,
		LoanID:      transactionScanner.LoanID.Int64,
		Amount:      transactionScanner.Amount.Float64,
		Description: transactionScanner.Description.String,
		CreatedAt:   transactionScanner.CreatedAt.Time,
	}

	return result, nil
}

func (r *transactionRepository) GetTransactionsByConsumerID(ctx context.Context, consumerID int64) (results []Transaction, err error) {
	query := `
		SELECT
			transaction_id,
			consumer_id,
			loan_id,
			amount,
			description,
			created_at
		FROM transactions
		WHERE deleted_at IS NULL
		AND consumer_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, consumerID)
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][GetTransactionsByConsumerID] while query. Err: %v", err))
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionScanner TransactionScanner
		err = rows.Scan(
			&transactionScanner.ID,
			&transactionScanner.ConsumerID,
			&transactionScanner.LoanID,
			&transactionScanner.Amount,
			&transactionScanner.Description,
			&transactionScanner.CreatedAt,
		)
		if err != nil {
			logger.Error(fmt.Sprintf("[transactionRepository][GetTransactionsByConsumerID] while scan query row. Err: %v", err))
			return results, err
		}

		result := Transaction{
			ID:          transactionScanner.ID.Int64,
			ConsumerID:  transactionScanner.ConsumerID.Int64,
			LoanID:      transactionScanner.LoanID.Int64,
			Amount:      transactionScanner.Amount.Float64,
			Description: transactionScanner.Description.String,
			CreatedAt:   transactionScanner.CreatedAt.Time,
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *transactionRepository) GetTransactionsByLoanID(ctx context.Context, loanID int64) (results []Transaction, err error) {
	query := `
		SELECT
			transaction_id,
			consumer_id,
			loan_id,
			amount,
			description,
			created_at
		FROM transactions
		WHERE deleted_at IS NULL
		AND loan_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, loanID)
	if err != nil {
		logger.Error(fmt.Sprintf("[transactionRepository][GetTransactionsByLoanID] while query. Err: %v", err))
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var transactionScanner TransactionScanner
		err = rows.Scan(
			&transactionScanner.ID,
			&transactionScanner.ConsumerID,
			&transactionScanner.LoanID,
			&transactionScanner.Amount,
			&transactionScanner.Description,
			&transactionScanner.CreatedAt,
		)
		if err != nil {
			logger.Error(fmt.Sprintf("[transactionRepository][GetTransactionsByLoanID] while scan query row. Err: %v", err))
			return results, err
		}

		result := Transaction{
			ID:          transactionScanner.ID.Int64,
			ConsumerID:  transactionScanner.ConsumerID.Int64,
			LoanID:      transactionScanner.LoanID.Int64,
			Amount:      transactionScanner.Amount.Float64,
			Description: transactionScanner.Description.String,
			CreatedAt:   transactionScanner.CreatedAt.Time,
		}

		results = append(results, result)
	}

	return results, nil
}
