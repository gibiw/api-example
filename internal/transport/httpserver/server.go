package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gihub.com/gibiw/api-example/internal/config"
	"gihub.com/gibiw/api-example/internal/entities"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type usecases interface {
	GetCars(ctx context.Context) ([]entities.Car, error)
	GetCarById(ctx context.Context, id uuid.UUID) (entities.Car, error)
	AddCar(ctx context.Context, car entities.Car) (entities.Car, error)
	DeleteCarById(ctx context.Context, id uuid.UUID) error
	UpdateCar(ctx context.Context, car entities.Car) (entities.Car, error)
}

type cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, error)
	Delete(key string)
}

type Server struct {
	cfg      config.Service
	usc      usecases
	ch       cache
	cacheTtl time.Duration
}

// TODO add tests and logs
func New(cfg config.Service, ucs usecases, ch cache) *Server {
	return &Server{
		cfg: cfg,
		usc: ucs,
		ch:  ch,
		cacheTtl: time.Second * time.Duration(cfg.CacheTtlSeconds),
	}
}

func (s *Server) Run() error {
	handlers := s.addHandlers()
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port), handlers)
}

func (s *Server) addHandlers() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(setResponseHeader())

	r.Route("/cars", func(r chi.Router) {
		r.Get("/", s.getCars())
		r.Post("/", s.addCar())
		r.Put("/", s.updateCar())

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", s.getCarById())
			r.Delete("/", s.deleteCarById())
		})
	})

	return r
}
