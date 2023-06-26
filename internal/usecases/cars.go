package usecases

import (
	"context"

	"gihub.com/gibiw/api-example/internal/entities"
	"github.com/google/uuid"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type repository interface {
	GetCars(ctx context.Context) ([]entities.Car, error)
	GetCarById(ctx context.Context, id uuid.UUID) (entities.Car, error)
	AddCar(ctx context.Context, car entities.Car) (entities.Car, error)
	DeleteCarById(ctx context.Context, id uuid.UUID) error
	UpdateCar(ctx context.Context, car entities.Car) (entities.Car, error)
}

// TODO add logs
type CarsUsecases struct {
	r repository
}

func New(r repository) *CarsUsecases {
	return &CarsUsecases{
		r: r,
	}
}

func (c *CarsUsecases) GetCars(ctx context.Context) ([]entities.Car, error) {
	return c.r.GetCars(ctx)
}

func (c *CarsUsecases) GetCarById(ctx context.Context, id uuid.UUID) (entities.Car, error) {
	return c.r.GetCarById(ctx, id)
}

func (c *CarsUsecases) AddCar(ctx context.Context, car entities.Car) (entities.Car, error) {
	return c.r.AddCar(ctx, car)
}

func (c *CarsUsecases) DeleteCarById(ctx context.Context, id uuid.UUID) error {
	return c.r.DeleteCarById(ctx, id)
}

func (c *CarsUsecases) UpdateCar(ctx context.Context, car entities.Car) (entities.Car, error) {
	return c.r.UpdateCar(ctx, car)
}
