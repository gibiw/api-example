package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gihub.com/gibiw/api-example/internal/entities"
	mycache "github.com/gibiw/cache"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gookit/slog"
)

func (s *Server) getCars() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cars, err := s.usc.GetCars(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		dtos := make([]CarDto, 0, len(cars))
		for _, v := range cars {
			dtos = append(dtos, carDomainToDto(v))
		}

		data, err := json.Marshal(dtos)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func (s *Server) getCarById() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		if value, ok := s.getValueFromCache(idParam); ok {
			car, err := json.Marshal(carDomainToDto(value))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			w.WriteHeader(http.StatusOK)
			w.Write(car)
			return
		}

		c, err := s.usc.GetCarById(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		s.ch.Set(idParam, c, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(c))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (s *Server) addCar() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		car := NewCarDto{}
		err = json.Unmarshal(body, &car)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		newCar, err := s.usc.AddCar(r.Context(), newCarToDomain(car))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		s.ch.Set(newCar.Id.String(), newCar, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(newCar))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}

func (s *Server) deleteCarById() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		err = s.usc.DeleteCarById(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) updateCar() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		defer r.Body.Close()

		car := CarDto{}
		err = json.Unmarshal(body, &car)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		newCar, err := s.usc.UpdateCar(r.Context(), carToDomain(car))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		s.ch.Set(newCar.Id.String(), newCar, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(newCar))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func (s *Server) getValueFromCache(key string) (entities.Car, bool) {
	value, err := s.ch.Get(key)
	if err != nil {
		if errors.Is(mycache.ErrorExpired, err) {
			slog.Debug("cache is expired for record with id ", key)
			s.ch.Delete(key)
			return entities.Car{}, false
		}

		slog.Debug(fmt.Sprintf("can not get record with id %s from cache: %s", key, err))

		return entities.Car{}, false
	}

	return value.(entities.Car), true
}
