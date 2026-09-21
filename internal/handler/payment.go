package handler

import (
	"net/http"

	"github.com/dutik/auto-hub/internal/dto"
	"github.com/dutik/auto-hub/internal/service"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	status := queryParam(r, "status")
	resp, err := h.svc.List(r.Context(), status)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	var req dto.PayPaymentRequest
	if err := DecodeAndValidate(r, &req); err != nil {
		HandleError(w, err)
		return
	}

	resp, err := h.svc.Pay(r.Context(), &req)
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}
