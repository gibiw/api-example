package usecases

import (
	"context"
	"errors"
	"testing"

	"gihub.com/gibiw/api-example/internal/entities"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCarsUsecases_GetCars(t *testing.T) {
	t.Run("get cars without error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		cars := []entities.Car{
			{
				Id:    uuid.New(),
				Brand: "Audi",
				Model: "A3",
				Color: "Red",
				Cost:  10000,
			},
			{
				Id:    uuid.New(),
				Brand: "Ford",
				Model: "Focus",
				Color: "Green",
				Cost:  8000,
			},
		}
		f.repository.EXPECT().GetCars(gomock.Any()).Return(cars, nil)
		usc := New(f.repository)

		// Act
		reps, err := usc.GetCars(context.Background())

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, cars, reps)
	})

	t.Run("get cars with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		returnErr := errors.New("text string")
		f.repository.EXPECT().GetCars(gomock.Any()).Return(nil, returnErr)
		usc := New(f.repository)

		// Act
		reps, err := usc.GetCars(context.Background())

		// Assert
		assert.Nil(t, reps)
		assert.Error(t, returnErr, err)
	})
}

func TestCarsUsecases_GetCarById(t *testing.T) {
	t.Run("get car by id without error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		id := uuid.New()
		car := entities.Car{

			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		f.repository.EXPECT().GetCarById(gomock.Any(), id).Return(car, nil)
		usc := New(f.repository)

		// Act
		reps, err := usc.GetCarById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, car, reps)
	})

	t.Run("get car by id with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		returnErr := errors.New("text string")
		id := uuid.New()
		car := entities.Car{}
		f.repository.EXPECT().GetCarById(gomock.Any(), id).Return(car, returnErr)
		usc := New(f.repository)

		// Act
		reps, err := usc.GetCarById(context.Background(), id)

		// Assert
		assert.Equal(t, car, reps)
		assert.Error(t, returnErr, err)
	})
}

func TestCarsUsecases_AddCar(t *testing.T) {
	t.Run("add car without error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		car := entities.Car{
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		f.repository.EXPECT().AddCar(gomock.Any(), car).Return(car, nil)
		usc := New(f.repository)

		// Act
		reps, err := usc.AddCar(context.Background(), car)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, car, reps)
	})

	t.Run("add car with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		returnErr := errors.New("text string")
		car := entities.Car{
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		f.repository.EXPECT().AddCar(gomock.Any(), car).Return(entities.Car{}, returnErr)
		usc := New(f.repository)

		// Act
		reps, err := usc.AddCar(context.Background(), car)

		// Assert
		assert.Equal(t, entities.Car{}, reps)
		assert.Error(t, returnErr, err)
	})
}

func TestCarsUsecases_DeleteCarById(t *testing.T) {
	t.Run("delete car by id without error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		id := uuid.New()
		f.repository.EXPECT().DeleteCarById(gomock.Any(), id).Return(nil)
		usc := New(f.repository)

		// Act
		err := usc.DeleteCarById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("delete car by id with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		returnErr := errors.New("text string")
		id := uuid.New()
		f.repository.EXPECT().DeleteCarById(gomock.Any(), id).Return(returnErr)
		usc := New(f.repository)

		// Act
		err := usc.DeleteCarById(context.Background(), id)

		// Assert
		assert.Error(t, returnErr, err)
	})
}

func TestCarsUsecases_UpdateCar(t *testing.T) {
	t.Run("update car without error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		id := uuid.New()
		car := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		f.repository.EXPECT().UpdateCar(gomock.Any(), car).Return(car, nil)
		usc := New(f.repository)

		// Act
		reps, err := usc.UpdateCar(context.Background(), car)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, car, reps)
	})

	t.Run("update car with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		returnErr := errors.New("text string")
		id := uuid.New()
		car := entities.Car{
			Id:    id,
			Brand: "Audi",
			Model: "A3",
			Color: "Red",
			Cost:  10000,
		}
		f.repository.EXPECT().UpdateCar(gomock.Any(), car).Return(entities.Car{}, returnErr)
		usc := New(f.repository)

		// Act
		reps, err := usc.UpdateCar(context.Background(), car)

		// Assert
		assert.Equal(t, entities.Car{}, reps)
		assert.Error(t, returnErr, err)
	})
}
