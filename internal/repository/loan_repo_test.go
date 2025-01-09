package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	repo := repository.NewLoanRepository(db)

	tests := []struct {
		name    string
		loan    repository.Loan
		wantID  int64
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			loan: repository.Loan{
				ConsumerLimitID: 1,
				ConsumerID:      1,
				MerchantID:      1,
				LoanAmount:      1000.0,
				ContractNumber:  "12345",
				InterestRate:    5.0,
				InterestAmount:  50.0,
				DueDate:         now,
				AssetName:       "Car",
			},
			wantID:  1,
			wantErr: false,
			mock: func() {
				mock.ExpectExec("INSERT INTO loans").
					WithArgs(1, 1, 1, 1000.0, "12345", 5.0, 50.0, sqlmock.AnyArg(), "Car").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "exec error",
			loan: repository.Loan{
				ConsumerLimitID: 1,
				ConsumerID:      1,
				MerchantID:      1,
				LoanAmount:      1000.0,
				ContractNumber:  "12345",
				InterestRate:    5.0,
				InterestAmount:  50.0,
				DueDate:         now,
				AssetName:       "Car",
			},
			wantID:  0,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("INSERT INTO loans").
					WithArgs(1, 1, 1, 1000.0, "12345", 5.0, 50.0, sqlmock.AnyArg(), "Car").
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name: "last insert id error",
			loan: repository.Loan{
				ConsumerLimitID: 1,
				ConsumerID:      1,
				MerchantID:      1,
				LoanAmount:      1000.0,
				ContractNumber:  "12345",
				InterestRate:    5.0,
				InterestAmount:  50.0,
				DueDate:         now,
				AssetName:       "Car",
			},
			wantID:  0,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("INSERT INTO loans").
					WithArgs(1, 1, 1, 1000.0, "12345", 5.0, 50.0, sqlmock.AnyArg(), "Car").
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			id, err := repo.CreateLoan(context.Background(), tt.loan)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantID, id)
		})
	}
}

func TestUpdateLoan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	trx, _ := db.Begin()

	repo := repository.NewLoanRepository(db)

	tests := []struct {
		name    string
		req     repository.UpdateLoanRequest
		tx      *sql.Tx
		wantErr bool
		mock    func()
	}{
		{
			name: "success with transaction",
			req: repository.UpdateLoanRequest{
				ID:                 1,
				PaidLoanAmount:     500.0,
				PaidInterestAmount: 25.0,
				LoanStatus:         "on_going",
				Installment:        5,
			},
			tx:      trx,
			wantErr: false,
			mock: func() {
				mock.ExpectExec("UPDATE loans").
					WithArgs(500.0, 25.0, "on_going", 5, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "success without transaction",
			req: repository.UpdateLoanRequest{
				ID:                 1,
				PaidLoanAmount:     500.0,
				PaidInterestAmount: 25.0,
				LoanStatus:         "on_going",
				Installment:        5,
			},
			tx:      nil,
			wantErr: false,
			mock: func() {
				mock.ExpectExec("UPDATE loans").
					WithArgs(500.0, 25.0, "on_going", 5, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "exec error with transaction",
			req: repository.UpdateLoanRequest{
				ID:                 1,
				PaidLoanAmount:     500.0,
				PaidInterestAmount: 25.0,
				LoanStatus:         "on_going",
				Installment:        5,
			},
			tx:      trx,
			wantErr: true,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE loans").
					WithArgs(500.0, 25.0, "on_going", 5, 1).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
		},
		{
			name: "exec error without transaction",
			req: repository.UpdateLoanRequest{
				ID:                 1,
				PaidLoanAmount:     500.0,
				PaidInterestAmount: 25.0,
				LoanStatus:         "on_going",
				Installment:        5,
			},
			tx:      nil,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("UPDATE loans").
					WithArgs(500.0, 25.0, "on_going", 5, 1).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdateLoan(context.Background(), tt.req, tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLoanByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewLoanRepository(db)
	now := time.Now()

	tests := []struct {
		name    string
		loanID  int64
		want    repository.Loan
		wantErr bool
		mock    func()
	}{
		{
			name:   "success",
			loanID: 1,
			want: repository.Loan{
				ID:                 1,
				ConsumerLimitID:    1,
				ConsumerID:         1,
				MerchantID:         1,
				LoanAmount:         1000.0,
				PaidLoanAmount:     500.0,
				ContractNumber:     "12345",
				InterestRate:       5.0,
				InterestAmount:     50.0,
				PaidInterestAmount: 25.0,
				LoanStatus:         "on_going",
				DueDate:            now,
				Installment:        5,
				AssetName:          "Car",
				CreatedAt:          now,
				UpdatedAt:          now,
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"loan_id", "consumer_limit_id", "consumer_id", "merchant_id", "loan_amount", "paid_loan_amount",
					"contract_number", "interest_rate", "interest_amount", "paid_interest_amount", "loan_status",
					"due_date", "installment", "asset_name", "created_at", "updated_at",
				}).AddRow(
					1, 1, 1, 1, 1000.0, 500.0, "12345", 5.0, 50.0, 25.0, "on_going",
					now, 5, "Car", now, now,
				)
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name:    "no rows",
			loanID:  2,
			want:    repository.Loan{},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(2).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:    "query error",
			loanID:  3,
			want:    repository.Loan{},
			wantErr: true,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(3).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetLoanByID(context.Background(), tt.loanID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteLoan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewLoanRepository(db)

	tests := []struct {
		name    string
		loanID  int64
		wantErr bool
		mock    func()
	}{
		{
			name:    "success",
			loanID:  1,
			wantErr: false,
			mock: func() {
				mock.ExpectExec("UPDATE loans").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "exec error",
			loanID:  2,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("UPDATE loans").
					WithArgs(2).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.DeleteLoan(context.Background(), tt.loanID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLoanByConsumerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewLoanRepository(db)
	now := time.Now()

	tests := []struct {
		name       string
		consumerID int64
		want       []repository.Loan
		wantErr    bool
		mock       func()
	}{
		{
			name:       "success",
			consumerID: 1,
			want: []repository.Loan{
				{
					ID:                 1,
					ConsumerLimitID:    1,
					ConsumerID:         1,
					MerchantID:         1,
					LoanAmount:         1000.0,
					PaidLoanAmount:     500.0,
					ContractNumber:     "12345",
					InterestRate:       5.0,
					InterestAmount:     50.0,
					PaidInterestAmount: 25.0,
					LoanStatus:         "on_going",
					DueDate:            now,
					Installment:        5,
					AssetName:          "Car",
					CreatedAt:          now,
					UpdatedAt:          now,
				},
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"loan_id", "consumer_limit_id", "consumer_id", "merchant_id", "loan_amount", "paid_loan_amount",
					"contract_number", "interest_rate", "interest_amount", "paid_interest_amount", "loan_status",
					"due_date", "installment", "asset_name", "created_at", "updated_at",
				}).AddRow(
					1, 1, 1, 1, 1000.0, 500.0, "12345", 5.0, 50.0, 25.0, "on_going",
					now, 5, "Car", now, now,
				)
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name:       "no rows",
			consumerID: 2,
			want:       nil,
			wantErr:    false,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(2).WillReturnRows(sqlmock.NewRows(nil))
			},
		},
		{
			name:       "query error",
			consumerID: 3,
			want:       nil,
			wantErr:    true,
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM loans WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(3).WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetLoanByConsumerID(context.Background(), tt.consumerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
