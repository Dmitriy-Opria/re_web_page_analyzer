package api

import (
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
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
	ctx := r.Context()
	request := ctx.Value(ContextUrlPayload).(*model.ParserRequest)
	resp, err := h.service.Parse(r.Context(), request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ServerErrorResponse(err, "can't execute service health")
	}
	return respondWithJson(w, http.StatusOK, resp)
}
