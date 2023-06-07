package entities

import "github.com/google/uuid"

type Car struct {
	Id    uuid.UUID `db:"id"`
	Brand string    `db:"brand"`
	Model string    `db:"model"`
	Color string    `db:"color"`
	Cost  uint64    `db:"cost"`
}
