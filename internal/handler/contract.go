package handler

import (
	"net/http"
	"strconv"

	"github.com/dutik/auto-hub/internal/dto"
	"github.com/dutik/auto-hub/internal/service"
	"github.com/go-chi/chi/v5"
)

type ContractHandler struct {
	svc *service.ContractService
}

func NewContractHandler(svc *service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

func (h *ContractHandler) List(w http.ResponseWriter, r *http.Request) {
	status := queryParam(r, "status")
	resp, err := h.svc.List(r.Context(), status)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateContractRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	contract, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusCreated, contract)
}

func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateContractRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	resp, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *ContractHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	payments, err := h.svc.GetPayments(r.Context(), id)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, payments)
}
