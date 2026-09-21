package handler

import (
	"net/http"
	"strconv"

	"github.com/dutik/auto-hub/internal/dto"
	"github.com/dutik/auto-hub/internal/service"
	"github.com/go-chi/chi/v5"
)

type FineHandler struct {
	svc *service.FineService
}

func NewFineHandler(svc *service.FineService) *FineHandler {
	return &FineHandler{svc: svc}
}

func (h *FineHandler) List(w http.ResponseWriter, r *http.Request) {
	status := queryParam(r, "status")

	var carID *int
	if v := r.URL.Query().Get("car_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			Error(w, http.StatusBadRequest, "invalid car_id")
			return
		}
		carID = &id
	}

	var driverID *int
	if v := r.URL.Query().Get("driver_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			Error(w, http.StatusBadRequest, "invalid driver_id")
			return
		}
		driverID = &id
	}

	resp, err := h.svc.List(r.Context(), status, carID, driverID)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *FineHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFineRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	fine, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, fine)
}

func (h *FineHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateFineRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	fine, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, fine)
}
