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

func TestCreateConsumer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerRepository(db)

	t.Run("success", func(t *testing.T) {
		consumer := repository.Consumer{
			FullName:     "John Doe",
			LegalName:    "Johnathan Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			Salary:       50000,
			NIK:          "1234567890",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
			CreatedAt:    time.Now(),
		}

		mock.ExpectExec("INSERT INTO consumers").
			WithArgs(
				consumer.FullName,
				consumer.LegalName,
				consumer.PlaceOfBirth,
				consumer.DOB,
				consumer.Salary,
				consumer.NIK,
				consumer.KTPImageURL,
				consumer.SelfieURL,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.CreateConsumer(context.Background(), consumer)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("failure", func(t *testing.T) {
		consumer := repository.Consumer{
			FullName:     "John Doe",
			LegalName:    "Johnathan Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			Salary:       50000,
			NIK:          "1234567890",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
			CreatedAt:    time.Now(),
		}

		mock.ExpectExec("INSERT INTO consumers").
			WithArgs(
				consumer.FullName,
				consumer.LegalName,
				consumer.PlaceOfBirth,
				consumer.DOB,
				consumer.Salary,
				consumer.NIK,
				consumer.KTPImageURL,
				consumer.SelfieURL,
			).
			WillReturnError(sql.ErrConnDone)

		id, err := repo.CreateConsumer(context.Background(), consumer)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestGetConsumerByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerRepository(db)

	t.Run("success", func(t *testing.T) {
		consumerID := int64(1)
		rows := sqlmock.NewRows([]string{
			"consumer_id", "full_name", "legal_name", "place_of_birth", "dob", "salary", "nik", "ktp_image_url", "selfie_image_url", "created_at",
		}).AddRow(
			consumerID, "John Doe", "Johnathan Doe", "New York", "1990-01-01", 50000, "1234567890", "http://example.com/ktp.jpg", "http://example.com/selfie.jpg", time.Now(),
		)

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null AND consumer_id = ?").
			WithArgs(consumerID).
			WillReturnRows(rows)

		consumer, err := repo.GetConsumerByID(context.Background(), consumerID)
		assert.NoError(t, err)
		assert.Equal(t, consumerID, consumer.ID)
		assert.Equal(t, "John Doe", consumer.FullName)
		assert.Equal(t, "Johnathan Doe", consumer.LegalName)
		assert.Equal(t, "New York", consumer.PlaceOfBirth)
		assert.Equal(t, "1990-01-01", consumer.DOB)
		assert.Equal(t, 50000.0, consumer.Salary)
		assert.Equal(t, "1234567890", consumer.NIK)
		assert.Equal(t, "http://example.com/ktp.jpg", consumer.KTPImageURL)
		assert.Equal(t, "http://example.com/selfie.jpg", consumer.SelfieURL)
	})

	t.Run("not found", func(t *testing.T) {
		consumerID := int64(2)

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null AND consumer_id = ?").
			WithArgs(consumerID).
			WillReturnError(sql.ErrNoRows)

		consumer, err := repo.GetConsumerByID(context.Background(), consumerID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), consumer.ID)
	})

	t.Run("query error", func(t *testing.T) {
		consumerID := int64(3)

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null AND consumer_id = ?").
			WithArgs(consumerID).
			WillReturnError(sql.ErrConnDone)

		consumer, err := repo.GetConsumerByID(context.Background(), consumerID)
		assert.Error(t, err)
		assert.Equal(t, int64(0), consumer.ID)
	})
}

func TestUpdateConsumer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerRepository(db)

	t.Run("success", func(t *testing.T) {
		consumer := repository.Consumer{
			ID:           1,
			FullName:     "John Doe",
			LegalName:    "Johnathan Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			Salary:       50000,
			NIK:          "1234567890",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
		}

		mock.ExpectExec("UPDATE consumers").
			WithArgs(
				consumer.FullName,
				consumer.LegalName,
				consumer.PlaceOfBirth,
				consumer.DOB,
				consumer.Salary,
				consumer.NIK,
				consumer.KTPImageURL,
				consumer.SelfieURL,
				consumer.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateConsumer(context.Background(), consumer)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		consumer := repository.Consumer{
			ID:           1,
			FullName:     "John Doe",
			LegalName:    "Johnathan Doe",
			PlaceOfBirth: "New York",
			DOB:          "1990-01-01",
			Salary:       50000,
			NIK:          "1234567890",
			KTPImageURL:  "http://example.com/ktp.jpg",
			SelfieURL:    "http://example.com/selfie.jpg",
		}

		mock.ExpectExec("UPDATE consumers").
			WithArgs(
				consumer.FullName,
				consumer.LegalName,
				consumer.PlaceOfBirth,
				consumer.DOB,
				consumer.Salary,
				consumer.NIK,
				consumer.KTPImageURL,
				consumer.SelfieURL,
				consumer.ID,
			).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateConsumer(context.Background(), consumer)
		assert.Error(t, err)
	})
}

func TestDeleteConsumer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerRepository(db)

	t.Run("success", func(t *testing.T) {
		consumerID := int64(1)

		mock.ExpectExec("UPDATE consumers").
			WithArgs(consumerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteConsumer(context.Background(), consumerID)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		consumerID := int64(1)

		mock.ExpectExec("UPDATE consumers").
			WithArgs(consumerID).
			WillReturnError(sql.ErrConnDone)

		err := repo.DeleteConsumer(context.Background(), consumerID)
		assert.Error(t, err)
	})
}

func TestFetchConsumer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewConsumerRepository(db)

	t.Run("success", func(t *testing.T) {
		req := repository.FetchConsumerRequest{
			Limit:  10,
			Offset: 0,
		}

		rows := sqlmock.NewRows([]string{
			"consumer_id", "full_name", "legal_name", "place_of_birth", "dob", "salary", "nik", "ktp_image_url", "selfie_image_url", "created_at",
		}).AddRow(
			1, "John Doe", "Johnathan Doe", "New York", "1990-01-01", 50000, "1234567890", "http://example.com/ktp.jpg", "http://example.com/selfie.jpg", time.Now(),
		).AddRow(
			2, "Jane Doe", "Janet Doe", "Los Angeles", "1992-02-02", 60000, "0987654321", "http://example.com/ktp2.jpg", "http://example.com/selfie2.jpg", time.Now(),
		)

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null LIMIT \\? OFFSET \\?").
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(rows)

		consumers, err := repo.FetchConsumer(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, consumers, 2)
		assert.Equal(t, int64(1), consumers[0].ID)
		assert.Equal(t, "John Doe", consumers[0].FullName)
		assert.Equal(t, "Johnathan Doe", consumers[0].LegalName)
		assert.Equal(t, "New York", consumers[0].PlaceOfBirth)
		assert.Equal(t, "1990-01-01", consumers[0].DOB)
		assert.Equal(t, 50000.0, consumers[0].Salary)
		assert.Equal(t, "1234567890", consumers[0].NIK)
		assert.Equal(t, "http://example.com/ktp.jpg", consumers[0].KTPImageURL)
		assert.Equal(t, "http://example.com/selfie.jpg", consumers[0].SelfieURL)

		assert.Equal(t, int64(2), consumers[1].ID)
		assert.Equal(t, "Jane Doe", consumers[1].FullName)
		assert.Equal(t, "Janet Doe", consumers[1].LegalName)
		assert.Equal(t, "Los Angeles", consumers[1].PlaceOfBirth)
		assert.Equal(t, "1992-02-02", consumers[1].DOB)
		assert.Equal(t, 60000.0, consumers[1].Salary)
		assert.Equal(t, "0987654321", consumers[1].NIK)
		assert.Equal(t, "http://example.com/ktp2.jpg", consumers[1].KTPImageURL)
		assert.Equal(t, "http://example.com/selfie2.jpg", consumers[1].SelfieURL)
	})

	t.Run("no rows", func(t *testing.T) {
		req := repository.FetchConsumerRequest{
			Limit:  10,
			Offset: 0,
		}

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null LIMIT \\? OFFSET \\?").
			WithArgs(req.Limit, req.Offset).
			WillReturnRows(sqlmock.NewRows(nil))

		consumers, err := repo.FetchConsumer(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, consumers, 0)
	})

	t.Run("query error", func(t *testing.T) {
		req := repository.FetchConsumerRequest{
			Limit:  10,
			Offset: 0,
		}

		mock.ExpectQuery("SELECT consumer_id, full_name, legal_name, place_of_birth, dob, salary, nik, ktp_image_url, selfie_image_url, created_at FROM consumers WHERE deleted_at is null LIMIT \\? OFFSET \\?").
			WithArgs(req.Limit, req.Offset).
			WillReturnError(sql.ErrConnDone)

		consumers, err := repo.FetchConsumer(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, consumers)
	})
}
