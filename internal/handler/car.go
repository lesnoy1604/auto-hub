package handler

import (
	"net/http"
	"strconv"

	"github.com/dutik/auto-hub/internal/dto"
	"github.com/dutik/auto-hub/internal/service"
	"github.com/go-chi/chi/v5"
)

type CarHandler struct {
	svc *service.CarService
}

func NewCarHandler(svc *service.CarService) *CarHandler {
	return &CarHandler{svc: svc}
}

func (h *CarHandler) List(w http.ResponseWriter, r *http.Request) {
	status := queryParam(r, "status")
	search := queryParam(r, "search")

	resp, err := h.svc.List(r.Context(), status, search)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *CarHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *CarHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCarRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	car, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, car)
}

func (h *CarHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateCarRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	car, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, car)
}

func queryParam(r *http.Request, key string) *string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return nil
	}
	return &v
}
