package httpserver

import "github.com/google/uuid"

type NewCarDto struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
	Color string `json:"color"`
	Cost  uint64 `json:"cost"`
}

type CarDto struct {
	Id    uuid.UUID `json:"id"`
	Brand string    `json:"brand"`
	Model string    `json:"model"`
	Color string    `json:"color"`
	Cost  uint64    `json:"cost"`
}
