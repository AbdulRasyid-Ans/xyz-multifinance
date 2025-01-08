package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateMerchant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewMerchantRepository(db)

	t.Run("success", func(t *testing.T) {
		merchant := repository.Merchant{
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
			CreatedAt:    time.Now(),
		}

		mock.ExpectExec("INSERT INTO merchants").
			WithArgs(merchant.MerchantName, merchant.MerchantType).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.CreateMerchant(context.Background(), merchant)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("failure", func(t *testing.T) {
		merchant := repository.Merchant{
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
			CreatedAt:    time.Now(),
		}

		mock.ExpectExec("INSERT INTO merchants").
			WithArgs(merchant.MerchantName, merchant.MerchantType).
			WillReturnError(errors.New("insert failed"))

		id, err := repo.CreateMerchant(context.Background(), merchant)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestGetMerchantByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewMerchantRepository(db)
	query := "SELECT merchant_id, merchant_name, merchant_type, created_at FROM merchants WHERE deleted_at IS NULL AND merchant_id = ?"

	t.Run("success", func(t *testing.T) {
		merchantID := int64(1)
		merchant := repository.Merchant{
			ID:           merchantID,
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
			CreatedAt:    time.Now(),
		}

		rows := sqlmock.NewRows([]string{"merchant_id", "merchant_name", "merchant_type", "created_at"}).
			AddRow(merchant.ID, merchant.MerchantName, merchant.MerchantType, merchant.CreatedAt)

		mock.ExpectQuery(query).
			WithArgs(merchantID).
			WillReturnRows(rows)

		result, err := repo.GetMerchantByID(context.Background(), merchantID)
		assert.NoError(t, err)
		assert.Equal(t, merchant, result)
	})

	t.Run("no rows", func(t *testing.T) {
		merchantID := int64(1)

		mock.ExpectQuery(query).
			WithArgs(merchantID).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetMerchantByID(context.Background(), merchantID)
		assert.NoError(t, err)
		assert.Equal(t, repository.Merchant{}, result)
	})

	t.Run("query error", func(t *testing.T) {
		merchantID := int64(1)

		mock.ExpectQuery(query).
			WithArgs(merchantID).
			WillReturnError(errors.New("query failed"))

		result, err := repo.GetMerchantByID(context.Background(), merchantID)
		assert.Error(t, err)
		assert.Equal(t, repository.Merchant{}, result)
	})
}

func TestFetchMerchant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewMerchantRepository(db)
	query := "SELECT merchant_id, merchant_name, merchant_type, created_at FROM merchants WHERE deleted_at IS NULL"

	t.Run("success", func(t *testing.T) {
		req := repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}

		merchant := repository.Merchant{
			ID:           1,
			MerchantName: "Test Merchant",
			MerchantType: "Retail",
			CreatedAt:    time.Now(),
		}

		rows := sqlmock.NewRows([]string{"merchant_id", "merchant_name", "merchant_type", "created_at"}).
			AddRow(merchant.ID, merchant.MerchantName, merchant.MerchantType, merchant.CreatedAt)

		mock.ExpectQuery(query).
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(rows)

		result, err := repo.FetchMerchant(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, merchant, result[0])
	})

	t.Run("no rows", func(t *testing.T) {
		req := repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}

		rows := sqlmock.NewRows([]string{"merchant_id", "merchant_name", "merchant_type", "created_at"})

		mock.ExpectQuery(query).
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(rows)

		result, err := repo.FetchMerchant(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("query error", func(t *testing.T) {
		req := repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}

		mock.ExpectQuery(query).
			WithArgs(req.Limit, req.Offset).
			WillReturnError(errors.New("query failed"))

		result, err := repo.FetchMerchant(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("scan error", func(t *testing.T) {
		req := repository.FetchMerchantRequest{
			Limit:  10,
			Offset: 0,
		}

		rows := sqlmock.NewRows([]string{"merchant_id", "merchant_name", "merchant_type", "created_at"}).
			AddRow("invalid_id", "Test Merchant", "Retail", time.Now())

		mock.ExpectQuery(query).
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(rows)

		result, err := repo.FetchMerchant(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateMerchant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewMerchantRepository(db)
	query := "UPDATE merchants"

	t.Run("success", func(t *testing.T) {
		merchant := repository.Merchant{
			ID:           1,
			MerchantName: "Updated Merchant",
			MerchantType: "Retail",
		}

		mock.ExpectExec(query).
			WithArgs(merchant.MerchantName, merchant.MerchantType, merchant.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateMerchant(context.Background(), merchant)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		merchant := repository.Merchant{
			ID:           1,
			MerchantName: "Updated Merchant",
			MerchantType: "Retail",
		}

		mock.ExpectExec(query).
			WithArgs(merchant.MerchantName, merchant.MerchantType, merchant.ID).
			WillReturnError(errors.New("update failed"))

		err := repo.UpdateMerchant(context.Background(), merchant)
		assert.Error(t, err)
	})
}

func TestDeleteMerchant(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewMerchantRepository(db)
	query := "UPDATE merchants"

	t.Run("success", func(t *testing.T) {
		merchantID := int64(1)

		mock.ExpectExec(query).
			WithArgs(merchantID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteMerchant(context.Background(), merchantID)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		merchantID := int64(1)

		mock.ExpectExec(query).
			WithArgs(merchantID).
			WillReturnError(errors.New("delete failed"))

		err := repo.DeleteMerchant(context.Background(), merchantID)
		assert.Error(t, err)
	})
}
