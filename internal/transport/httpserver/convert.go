package httpserver

import "gihub.com/gibiw/api-example/internal/entities"

func newCarToDomain(nc NewCarDto) entities.Car {
	return entities.Car{
		Brand: nc.Brand,
		Model: nc.Model,
		Color: nc.Color,
		Cost:  nc.Cost,
	}
}

func carDomainToDto(c entities.Car) CarDto {
	return CarDto{
		Id:    c.Id,
		Brand: c.Brand,
		Model: c.Model,
		Color: c.Color,
		Cost:  c.Cost,
	}
}

func carToDomain(c CarDto) entities.Car {
	return entities.Car{
		Id:    c.Id,
		Brand: c.Brand,
		Model: c.Model,
		Color: c.Color,
		Cost:  c.Cost,
	}
}
