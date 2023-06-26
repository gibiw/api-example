package repository

import (
	"context"

	"gihub.com/gibiw/api-example/internal/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	getAllCarsQuery = "SELECT id, brand, model, color, cost FROM cars"
	getCarQuery     = "SELECT id, brand, model, color, cost FROM cars WHERE id=$1"
	addCarQuery     = "INSERT INTO cars (brand, model, color, cost) VALUES ($1, $2, $3, $4) RETURNING id, brand, model, color, cost"
	deleteCarQuery  = "DELETE FROM cars WHERE id=$1"
	updateCarQuery  = "UPDATE cars SET brand=$1, model=$2, color=$3, cost=$4 WHERE id=$5"
)

type CarRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) GetCars(ctx context.Context) ([]entities.Car, error) {
	cars := []entities.Car{}
	if err := r.db.SelectContext(ctx, &cars, getAllCarsQuery); err != nil {
		return nil, err
	}

	return cars, nil
}

func (r *CarRepository) GetCarById(ctx context.Context, id uuid.UUID) (entities.Car, error) {
	car := entities.Car{}

	if err := r.db.GetContext(ctx, &car, getCarQuery, id); err != nil {
		return entities.Car{}, err
	}

	return car, nil
}

func (r *CarRepository) AddCar(ctx context.Context, car entities.Car) (entities.Car, error) {
	newCar := entities.Car{}

	err := r.db.QueryRowxContext(ctx, addCarQuery, car.Brand, car.Model, car.Color, car.Cost).StructScan(&newCar)

	if err != nil {
		return entities.Car{}, err
	}

	return newCar, nil
}

func (r *CarRepository) DeleteCarById(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteCarQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *CarRepository) UpdateCar(ctx context.Context, car entities.Car) (entities.Car, error) {
	_, err := r.GetCarById(ctx, car.Id)
	if err != nil {
		return entities.Car{}, err
	}

	_, err = r.db.ExecContext(ctx, updateCarQuery, car.Brand, car.Model, car.Color, car.Cost, car.Id)

	if err != nil {
		return entities.Car{}, err
	}

	return car, nil
}
