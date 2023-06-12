package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"gihub.com/gibiw/api-example/internal/entities"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCarRepository_GetCars(t *testing.T) {
	t.Run("with rows", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectedCars := []entities.Car{
			{
				Id:    uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c"),
				Brand: "Audi",
				Model: "A3",
				Color: "Red",
				Cost:  10000,
			},
			{
				Id:    uuid.MustParse("3d997272-468f-4b66-91db-00c39f0ef717"),
				Brand: "BMW",
				Model: "X6",
				Color: "Black",
				Cost:  20000,
			},
		}
		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"}).
			AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "Audi", "A3", "Red", 10000).
			AddRow("3d997272-468f-4b66-91db-00c39f0ef717", "BMW", "X6", "Black", 20000)

		f.mock.ExpectQuery("SELECT id, brand, model, color, cost FROM cars").
			WillReturnRows(rows)
		repo := New(f.db)

		// Act
		cars, err := repo.GetCars(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedCars, cars)
	})

	t.Run("without rows", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"})

		f.mock.ExpectQuery("SELECT id, brand, model, color, cost FROM cars").
			WillReturnRows(rows)
		repo := New(f.db)

		// Act
		cars, err := repo.GetCars(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, []entities.Car{}, cars)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		f.mock.ExpectQuery("SELECT id, brand, model, color, cost FROM cars").
			WillReturnError(expectErr)
		repo := New(f.db)

		// Act
		cars, err := repo.GetCars(context.Background())

		// Assert
		assert.Error(t, expectErr, err)
		assert.ElementsMatch(t, []entities.Car{}, cars)
	})
}

func TestCarRepository_GetCarById(t *testing.T) {
	t.Run("with car", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")
		expectedCar := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"}).
			AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "Audi", "A3", "Red", 10000)

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.db)

		// Act
		car, err := repo.GetCarById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCar, car)
	})

	t.Run("without car", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")

		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"})

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.db)

		// Act
		car, err := repo.GetCarById(context.Background(), id)

		// Assert
		assert.Error(t, sql.ErrNoRows, err)
		assert.Equal(t, entities.Car{}, car)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.db)

		// Act
		car, err := repo.GetCarById(context.Background(), id)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, entities.Car{}, car)
	})
}

func TestCarRepository_AddCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		expectedCar := entities.Car{
			Id:    uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c"),
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"}).
			AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "Audi", "A3", "Red", 10000)

		f.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO cars (brand, model, color, cost) VALUES ($1, $2, $3, $4) RETURNING id, brand, model, color, cost")).
			WithArgs(expectedCar.Brand, expectedCar.Model, expectedCar.Color, expectedCar.Cost).
			WillReturnRows(rows)

		repo := New(f.db)

		// Act
		car, err := repo.AddCar(context.Background(), expectedCar)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCar, car)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		expectedCar := entities.Car{
			Id:    uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c"),
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}

		f.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO cars (brand, model, color, cost) VALUES ($1, $2, $3, $4) RETURNING id, brand, model, color, cost")).
			WithArgs(expectedCar.Brand, expectedCar.Model, expectedCar.Color, expectedCar.Cost).
			WillReturnError(expectErr)
		repo := New(f.db)

		// Act
		car, err := repo.AddCar(context.Background(), expectedCar)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, entities.Car{}, car)
	})
}

func TestCarRepository_DeleteCarById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")

		f.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(f.db)

		// Act
		err := repo.DeleteCarById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")

		f.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)

		repo := New(f.db)

		// Act
		err := repo.DeleteCarById(context.Background(), id)

		// Assert
		assert.ErrorIs(t, expectErr, err)
	})
}

func TestCarRepository_UpdateCar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")
		expectedCar := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"}).
			AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "Audi", "A3", "Red", 10000)

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)

		f.mock.ExpectExec(regexp.QuoteMeta("UPDATE cars SET brand=$1, model=$2, color=$3, cost=$4 WHERE id=$5")).
			WithArgs(expectedCar.Brand, expectedCar.Model, expectedCar.Color, expectedCar.Cost, expectedCar.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(f.db)

		// Act
		car, err := repo.UpdateCar(context.Background(), expectedCar)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCar, car)
	})

	t.Run("without car", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")
		expectedCar := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}

		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"})

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).WithArgs(id).WillReturnRows(rows)
		repo := New(f.db)

		// Act
		car, err := repo.UpdateCar(context.Background(), expectedCar)

		// Assert
		assert.Error(t, sql.ErrNoRows, err)
		assert.Equal(t, entities.Car{}, car)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		id := uuid.MustParse("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c")
		expectedCar := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		rows := sqlmock.NewRows([]string{"id", "brand", "model", "color", "cost"}).
			AddRow("bea1b24d-0627-4ea0-aa2b-8af4c6c2a41c", "Audi", "A3", "Red", 10000)

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, brand, model, color, cost FROM cars WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)

		f.mock.ExpectExec(regexp.QuoteMeta("UPDATE cars SET brand=$1, model=$2, color=$3, cost=$4 WHERE id=$5")).
			WithArgs(expectedCar.Brand, expectedCar.Model, expectedCar.Color, expectedCar.Cost, expectedCar.Id).
			WillReturnError(expectErr)

		repo := New(f.db)

		// Act
		car, err := repo.UpdateCar(context.Background(), expectedCar)

		// Assert
		assert.ErrorIs(t, expectErr, err)
		assert.Equal(t, entities.Car{}, car)
	})
}
