package handler

import (
	"net/http"
	"strconv"

	"github.com/dutik/auto-hub/internal/dto"
	"github.com/dutik/auto-hub/internal/service"
	"github.com/go-chi/chi/v5"
)

type DriverHandler struct {
	svc *service.DriverService
}

func NewDriverHandler(svc *service.DriverService) *DriverHandler {
	return &DriverHandler{svc: svc}
}

func (h *DriverHandler) List(w http.ResponseWriter, r *http.Request) {
	status := queryParam(r, "status")
	search := queryParam(r, "search")

	resp, err := h.svc.List(r.Context(), status, search)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *DriverHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

func (h *DriverHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDriverRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	driver, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, driver)
}

func (h *DriverHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateDriverRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	driver, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, driver)
}
