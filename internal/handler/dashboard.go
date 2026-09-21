package handler

import (
	"net/http"

	"github.com/dutik/auto-hub/internal/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.Get(r.Context())
	if err != nil {
		HandleError(w, err)
		return
	}
	JSON(w, http.StatusOK, resp)
}
