package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/log"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
	"github.com/pkg/errors"
)

type StaffHandler struct {
	service *service.StaffService
}

func NewStaffHandler(staffService *service.StaffService) *StaffHandler {
	return &StaffHandler{service: staffService}
}

func (h *StaffHandler) Health(w http.ResponseWriter, r *http.Request) *rye.Response {
	err := h.service.Health(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ServerErrorResponse(err, "can't execute service health")
	}
	return respondWithMessage(w, http.StatusOK, "health")
}

func respondWithJsonAndStatusCode(w http.ResponseWriter, statusCode int, payload interface{}) *rye.Response {
	return respondWithJson(w, statusCode, payload)
}

func respondWithMessage(w http.ResponseWriter, statusCode int, message string) *rye.Response {
	responseJSON := &rye.JSONStatus{
		Message: message,
		Status:  "OK",
	}
	return respondWithJson(w, statusCode, responseJSON)
}

func respondWithJson(w http.ResponseWriter, statusCode int, jsonData interface{}) *rye.Response {
	responseJSON, err := json.MarshalIndent(jsonData, "", "\t")
	if err != nil {
		return ServerErrorResponse(err, "unable to generate response JSON")
	}
	rye.WriteJSONResponse(w, statusCode, responseJSON)
	return nil
}

func ErrorResponse(err error, message string, httpStatusCode int) *rye.Response {
	if err == nil {
		err = fmt.Errorf(message)
	}
	return &rye.Response{
		Err:        errors.Wrap(err, message),
		StatusCode: httpStatusCode,
	}
}

func ServerErrorResponse(err error, message string) *rye.Response {
	if err != nil {
		log.Errorln(errors.Wrap(err, message))
	}
	return ErrorResponse(err, message, http.StatusInternalServerError)
}

func BadRequestResponse(err error, message string) *rye.Response {
	if err != nil {
		log.Errorln(errors.Wrap(err, message))
	}
	return WarningResponse(err, message, http.StatusBadRequest)
}

func NotFoundResponse(err error, message string) *rye.Response {
	if err != nil {
		log.Errorln(errors.Wrap(err, message))
	}
	return WarningResponse(err, message, http.StatusNotFound)
}

func WarningResponse(err error, message string, httpStatusCode int) *rye.Response {
	if err == nil {
		err = fmt.Errorf(message)
	} else if message != "" {
		err = errors.Wrap(err, message)
	}
	return &rye.Response{
		Err:        err,
		StatusCode: httpStatusCode,
	}
}
