package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLimitByTenureAndConsumerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)

	tests := []struct {
		name       string
		tenure     int16
		consumerID int64
		mock       func()
		want       repository.ConsumerLimit
		wantErr    bool
	}{
		{
			name:       "success",
			tenure:     1,
			consumerID: 123,
			mock: func() {
				rows := sqlmock.NewRows([]string{"consumer_limit_id", "consumer_id", "tenure", "limit_amount"}).
					AddRow(1, 123, 1, 1000.0)
				mock.ExpectQuery("SELECT consumer_limit_id, consumer_id, tenure, limit_amount FROM consumer_limits WHERE deleted_at IS NULL AND consumer_id = \\? AND tenure = \\? LIMIT 1").
					WithArgs(123, 1).
					WillReturnRows(rows)
			},
			want: repository.ConsumerLimit{
				ID:          1,
				ConsumerID:  123,
				Tenure:      1,
				LimitAmount: 1000.0,
			},
			wantErr: false,
		},
		{
			name:       "no rows",
			tenure:     1,
			consumerID: 123,
			mock: func() {
				mock.ExpectQuery("SELECT consumer_limit_id, consumer_id, tenure, limit_amount FROM consumer_limits WHERE deleted_at IS NULL AND consumer_id = \\? AND tenure = \\? LIMIT 1").
					WithArgs(123, 1).
					WillReturnError(sql.ErrNoRows)
			},
			want:    repository.ConsumerLimit{},
			wantErr: false,
		},
		{
			name:       "query error",
			tenure:     1,
			consumerID: 123,
			mock: func() {
				mock.ExpectQuery("SELECT consumer_limit_id, consumer_id, tenure, limit_amount FROM consumer_limits WHERE deleted_at IS NULL AND consumer_id = \\? AND tenure = \\? LIMIT 1").
					WithArgs(123, 1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    repository.ConsumerLimit{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetLimitByTenureAndConsumerID(context.Background(), tt.tenure, tt.consumerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateConsumerLimit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)
	query := "INSERT INTO consumer_limits"

	tests := []struct {
		name          string
		consumerLimit repository.ConsumerLimit
		mock          func()
		wantID        int64
		wantErr       bool
	}{
		{
			name: "success",
			consumerLimit: repository.ConsumerLimit{
				ConsumerID:  123,
				Tenure:      1,
				LimitAmount: 1000.0,
			},
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(123, 1, 1000.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantID:  1,
			wantErr: false,
		},
		{
			name: "insert error",
			consumerLimit: repository.ConsumerLimit{
				ConsumerID:  123,
				Tenure:      1,
				LimitAmount: 1000.0,
			},
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(123, 1, 1000.0).
					WillReturnError(sql.ErrConnDone)
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			gotID, err := repo.CreateConsumerLimit(context.Background(), tt.consumerLimit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantID, gotID)
		})
	}
}

func TestUpdateConsumerLimit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)
	query := "UPDATE consumer_limits"

	tests := []struct {
		name          string
		consumerLimit repository.ConsumerLimit
		mock          func()
		wantErr       bool
	}{
		{
			name: "success",
			consumerLimit: repository.ConsumerLimit{
				ID:          1,
				LimitAmount: 2000.0,
			},
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(2000.0, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "update error",
			consumerLimit: repository.ConsumerLimit{
				ID:          1,
				LimitAmount: 2000.0,
			},
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(2000.0, 1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdateConsumerLimit(context.Background(), tt.consumerLimit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteConsumerLimit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)
	query := "UPDATE consumer_limits"

	tests := []struct {
		name            string
		consumerLimitID int64
		mock            func()
		wantErr         bool
	}{
		{
			name:            "success",
			consumerLimitID: 1,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:            "delete error",
			consumerLimitID: 1,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.DeleteConsumerLimit(context.Background(), tt.consumerLimitID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetConsumerLimitByConsumerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)
	query := "SELECT consumer_limit_id, consumer_id, tenure, limit_amount FROM consumer_limits WHERE deleted_at IS NULL AND consumer_id = ?"

	tests := []struct {
		name       string
		consumerID int64
		mock       func()
		want       []repository.ConsumerLimit
		wantErr    bool
	}{
		{
			name:       "success",
			consumerID: 123,
			mock: func() {
				rows := sqlmock.NewRows([]string{"consumer_limit_id", "consumer_id", "tenure", "limit_amount"}).
					AddRow(1, 123, 1, 1000.0).
					AddRow(2, 123, 2, 2000.0)
				mock.ExpectQuery(query).
					WithArgs(123).
					WillReturnRows(rows)
			},
			want: []repository.ConsumerLimit{
				{
					ID:          1,
					ConsumerID:  123,
					Tenure:      1,
					LimitAmount: 1000.0,
				},
				{
					ID:          2,
					ConsumerID:  123,
					Tenure:      2,
					LimitAmount: 2000.0,
				},
			},
			wantErr: false,
		},
		{
			name:       "no rows",
			consumerID: 123,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs(123).WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:       "query error",
			consumerID: 123,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs(123).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetConsumerLimitByConsumerID(context.Background(), tt.consumerID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetConsumerLimitByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerLimitRepository(db)
	query := "SELECT consumer_limit_id, consumer_id, tenure, limit_amount FROM consumer_limits WHERE deleted_at IS NULL AND consumer_limit_id = ?"

	tests := []struct {
		name            string
		consumerLimitID int64
		mock            func()
		want            repository.ConsumerLimit
		wantErr         bool
	}{
		{
			name:            "success",
			consumerLimitID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"consumer_limit_id", "consumer_id", "tenure", "limit_amount"}).
					AddRow(1, 123, 1, 1000.0)
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: repository.ConsumerLimit{
				ID:          1,
				ConsumerID:  123,
				Tenure:      1,
				LimitAmount: 1000.0,
			},
			wantErr: false,
		},
		{
			name:            "no rows",
			consumerLimitID: 1,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			want:    repository.ConsumerLimit{},
			wantErr: false,
		},
		{
			name:            "query error",
			consumerLimitID: 1,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			want:    repository.ConsumerLimit{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetConsumerLimitByID(context.Background(), tt.consumerLimitID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
