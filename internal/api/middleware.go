package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
	"github.com/InVisionApp/rye"
)

type contextKey int

const (
	notUsed contextKey = iota

	ContextUrlPayload
)

func middlewareParsePayload(w http.ResponseWriter, r *http.Request) *rye.Response {
	ctx := r.Context()
	request := model.ParserRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return BadRequestResponse(err, "Error during reading request body")
	}
	if err := json.Unmarshal(body, &request); err != nil {
		return BadRequestResponse(err, "Error during unmarshalling request body")
	}
	ctx = context.WithValue(ctx, ContextUrlPayload, &request)
	return &rye.Response{Context: ctx}
}
