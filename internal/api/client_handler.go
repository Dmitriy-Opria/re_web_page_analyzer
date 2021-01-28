package api

import (
	"net/http"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
)

type ClientHandler struct {
	service *service.ClientService
}

func NewClientHandler(parseService *service.ClientService) *ClientHandler {
	return &ClientHandler{service: parseService}
}

func (h *ClientHandler) Client(w http.ResponseWriter, r *http.Request) *rye.Response {
	err := h.service.Get(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ServerErrorResponse(err, "can't execute service health")
	}
	return respondWithMessage(w, http.StatusOK, "health")
}
