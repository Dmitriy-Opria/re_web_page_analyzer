package api

import (
	"net/http"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
)

type ParserHandler struct {
	service *service.ParserService
}

func NewParserHandler(parseService *service.ParserService) *ParserHandler {
	return &ParserHandler{service: parseService}
}

func (h *ParserHandler) Parse(w http.ResponseWriter, r *http.Request) *rye.Response {
	err := h.service.Parse(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ServerErrorResponse(err, "can't execute service health")
	}
	return respondWithMessage(w, http.StatusOK, "health")
}
