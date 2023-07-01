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

// getCars godoc
// @Summary      Get all cars
// @Description  Get all cars
// @Tags         cars
// @Accept       json
// @Produce      json
// @Success      200  {object}  []CarDto
// @Failure      500  {object}  errorResponse
// @Router       /cars/ [get]
func (s *Server) getCars() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cars, err := s.usc.GetCars(r.Context())
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		dtos := make([]CarDto, 0, len(cars))
		for _, v := range cars {
			dtos = append(dtos, carDomainToDto(v))
		}

		data, err := json.Marshal(dtos)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// getCarById godoc
// @Summary      Get a car by ID
// @Description  Get a car by ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Car ID"
// @Success      200  {object}  CarDto
// @Failure      500  {object}  errorResponse
// @Router       /cars/{id} [get]
func (s *Server) getCarById() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		if value, ok := s.getValueFromCache(idParam); ok {
			car, err := json.Marshal(carDomainToDto(value))
			if err != nil {
				newErrorResponse(w, http.StatusInternalServerError, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(car)
			return
		}

		c, err := s.usc.GetCarById(r.Context(), id)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		s.ch.Set(idParam, c, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(c))
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// addCar godoc
// @Summary      Add new car
// @Description  Add new car
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        request    body      NewCarDto  true  "Car"
// @Success      201  {object}  CarDto
// @Failure      500  {object}  errorResponse
// @Router       /cars [post]
func (s *Server) addCar() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		car := NewCarDto{}
		err = json.Unmarshal(body, &car)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		newCar, err := s.usc.AddCar(r.Context(), newCarToDomain(car))
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		s.ch.Set(newCar.Id.String(), newCar, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(newCar))
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}

// deleteCarById godoc
// @Summary      Delete a car by ID
// @Description  Delete a car by ID
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Car ID"
// @Success      200
// @Failure      500  {object}  errorResponse
// @Router       /cars/{id} [delete]
func (s *Server) deleteCarById() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err = s.usc.DeleteCarById(r.Context(), id)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// updateCar godoc
// @Summary      Update a car
// @Description  Update a car
// @Tags         cars
// @Accept       json
// @Produce      json
// @Param        request    body      CarDto  true  "Car"
// @Success      200  {object}  CarDto
// @Failure      500  {object}  errorResponse
// @Router       /cars [put]
func (s *Server) updateCar() func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			newErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		defer r.Body.Close()

		car := CarDto{}
		err = json.Unmarshal(body, &car)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		newCar, err := s.usc.UpdateCar(r.Context(), carToDomain(car))
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		s.ch.Set(newCar.Id.String(), newCar, s.cacheTtl)

		resp, err := json.Marshal(carDomainToDto(newCar))
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, err)
			return
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
