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

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	trx, _ := db.Begin()

	repo := repository.NewTransactionRepository(db)

	tests := []struct {
		name    string
		tx      *sql.Tx
		input   repository.Transaction
		wantID  int64
		wantErr bool
		mock    func()
	}{
		{
			name: "success with transaction",
			tx:   trx,
			input: repository.Transaction{
				ConsumerID:  1,
				LoanID:      1,
				Amount:      1000.0,
				Description: "Test transaction",
			},
			wantID:  1,
			wantErr: false,
			mock: func() {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(1, 1, 1000.0, "Test transaction").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "exec error with transaction",
			tx:   trx,
			input: repository.Transaction{
				ConsumerID:  1,
				LoanID:      1,
				Amount:      1000.0,
				Description: "Test transaction",
			},
			wantID:  0,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(1, 1, 1000.0, "Test transaction").
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name: "last insert id error",
			tx:   trx,
			input: repository.Transaction{
				ConsumerID:  1,
				LoanID:      1,
				Amount:      1000.0,
				Description: "Test transaction",
			},
			wantID:  0,
			wantErr: true,
			mock: func() {
				mock.ExpectExec("INSERT INTO transactions").
					WithArgs(1, 1, 1000.0, "Test transaction").
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			id, err := repo.CreateTransaction(context.Background(), tt.tx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestGetTransactionByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	tests := []struct {
		name    string
		id      int64
		want    repository.Transaction
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			id:   1,
			want: repository.Transaction{
				ID:          1,
				ConsumerID:  1,
				LoanID:      1,
				Amount:      1000.0,
				Description: "Test transaction",
				CreatedAt:   time.Now(),
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"transaction_id", "consumer_id", "loan_id", "amount", "description", "created_at"}).
					AddRow(1, 1, 1, 1000.0, "Test transaction", time.Now())
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND transaction_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:    "no rows",
			id:      2,
			want:    repository.Transaction{},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND transaction_id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:    "query error",
			id:      3,
			want:    repository.Transaction{},
			wantErr: true,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND transaction_id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := repo.GetTransactionByID(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, result.ID)
				assert.Equal(t, tt.want.ConsumerID, result.ConsumerID)
				assert.Equal(t, tt.want.LoanID, result.LoanID)
				assert.Equal(t, tt.want.Amount, result.Amount)
				assert.Equal(t, tt.want.Description, result.Description)
				// Use assert.WithinDuration to compare time fields
				assert.WithinDuration(t, tt.want.CreatedAt, result.CreatedAt, time.Second)
			}
		})
	}
}

func TestGetTransactionsByConsumerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	tests := []struct {
		name       string
		consumerID int64
		want       []repository.Transaction
		wantErr    bool
		mock       func()
	}{
		{
			name:       "success",
			consumerID: 1,
			want: []repository.Transaction{
				{
					ID:          1,
					ConsumerID:  1,
					LoanID:      1,
					Amount:      1000.0,
					Description: "Test transaction 1",
					CreatedAt:   time.Now(),
				},
				{
					ID:          2,
					ConsumerID:  1,
					LoanID:      2,
					Amount:      2000.0,
					Description: "Test transaction 2",
					CreatedAt:   time.Now(),
				},
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"transaction_id", "consumer_id", "loan_id", "amount", "description", "created_at"}).
					AddRow(1, 1, 1, 1000.0, "Test transaction 1", time.Now()).
					AddRow(2, 1, 2, 2000.0, "Test transaction 2", time.Now())
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:       "no rows",
			consumerID: 2,
			want:       []repository.Transaction{},
			wantErr:    false,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(nil))
			},
		},
		{
			name:       "query error",
			consumerID: 3,
			want:       []repository.Transaction{},
			wantErr:    true,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:       "scan error",
			consumerID: 4,
			want:       []repository.Transaction{},
			wantErr:    true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"transaction_id", "consumer_id", "loan_id", "amount", "description", "created_at"}).
					AddRow("invalid", 4, 1, 1000.0, "Test transaction", time.Now())
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND consumer_id = ?").
					WithArgs(4).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			results, err := repo.GetTransactionsByConsumerID(context.Background(), tt.consumerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(results))
				for i, want := range tt.want {
					assert.Equal(t, want.ID, results[i].ID)
					assert.Equal(t, want.ConsumerID, results[i].ConsumerID)
					assert.Equal(t, want.LoanID, results[i].LoanID)
					assert.Equal(t, want.Amount, results[i].Amount)
					assert.Equal(t, want.Description, results[i].Description)
					assert.WithinDuration(t, want.CreatedAt, results[i].CreatedAt, time.Second)
				}
			}
		})
	}
}

func TestGetTransactionsByLoanID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	tests := []struct {
		name    string
		loanID  int64
		want    []repository.Transaction
		wantErr bool
		mock    func()
	}{
		{
			name:   "success",
			loanID: 1,
			want: []repository.Transaction{
				{
					ID:          1,
					ConsumerID:  1,
					LoanID:      1,
					Amount:      1000.0,
					Description: "Test transaction 1",
					CreatedAt:   time.Now(),
				},
				{
					ID:          2,
					ConsumerID:  2,
					LoanID:      1,
					Amount:      2000.0,
					Description: "Test transaction 2",
					CreatedAt:   time.Now(),
				},
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"transaction_id", "consumer_id", "loan_id", "amount", "description", "created_at"}).
					AddRow(1, 1, 1, 1000.0, "Test transaction 1", time.Now()).
					AddRow(2, 2, 1, 2000.0, "Test transaction 2", time.Now())
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:    "no rows",
			loanID:  2,
			want:    []repository.Transaction{},
			wantErr: false,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows(nil))
			},
		},
		{
			name:    "query error",
			loanID:  3,
			want:    []repository.Transaction{},
			wantErr: true,
			mock: func() {
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(3).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:    "scan error",
			loanID:  4,
			want:    []repository.Transaction{},
			wantErr: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"transaction_id", "consumer_id", "loan_id", "amount", "description", "created_at"}).
					AddRow("invalid", 4, 1, 1000.0, "Test transaction", time.Now())
				mock.ExpectQuery("SELECT transaction_id, consumer_id, loan_id, amount, description, created_at FROM transactions WHERE deleted_at IS NULL AND loan_id = ?").
					WithArgs(4).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			results, err := repo.GetTransactionsByLoanID(context.Background(), tt.loanID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(results))
				for i, want := range tt.want {
					assert.Equal(t, want.ID, results[i].ID)
					assert.Equal(t, want.ConsumerID, results[i].ConsumerID)
					assert.Equal(t, want.LoanID, results[i].LoanID)
					assert.Equal(t, want.Amount, results[i].Amount)
					assert.Equal(t, want.Description, results[i].Description)
					assert.WithinDuration(t, want.CreatedAt, results[i].CreatedAt, time.Second)
				}
			}
		})
	}
}

func TestBeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	tests := []struct {
		name    string
		wantErr bool
		mock    func()
	}{
		{
			name:    "success",
			wantErr: false,
			mock: func() {
				mock.ExpectBegin()
			},
		},
		{
			name:    "begin error",
			wantErr: true,
			mock: func() {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			tx, err := repo.BeginTx(context.Background())
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
			}
		})
	}
}

func TestCommitTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	mock.ExpectBegin()
	trx, _ := db.Begin()

	tests := []struct {
		name    string
		tx      *sql.Tx
		wantErr bool
		mock    func()
	}{
		{
			name:    "success",
			tx:      trx,
			wantErr: false,
			mock: func() {
				mock.ExpectCommit()
			},
		},
		{
			name:    "commit error",
			tx:      trx,
			wantErr: true,
			mock: func() {
				mock.ExpectCommit().WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.CommitTx(context.Background(), tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRollbackTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewTransactionRepository(db)

	mock.ExpectBegin()
	trx, _ := db.Begin()

	tests := []struct {
		name    string
		tx      *sql.Tx
		wantErr bool
		mock    func()
	}{
		{
			name:    "success",
			tx:      trx,
			wantErr: false,
			mock: func() {
				mock.ExpectRollback()
			},
		},
		{
			name:    "rollback error",
			tx:      trx,
			wantErr: true,
			mock: func() {
				mock.ExpectRollback().WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.RollbackTx(context.Background(), tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
