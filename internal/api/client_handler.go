package api

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
)

type ClientHandler struct {
	parser *service.ParserService
}

func NewClientHandler(parser *service.ParserService) *ClientHandler {
	return &ClientHandler{
		parser: parser,
	}
}

func (h *ClientHandler) Form(w http.ResponseWriter, r *http.Request) *rye.Response {
	tmpl := template.Must(template.ParseFiles("../static/form.html"))
	tmpl.Execute(w, nil)
	return nil
}

func (h *ClientHandler) Result(w http.ResponseWriter, r *http.Request) *rye.Response {
	url := r.FormValue("url")
	if len(url) == 0 {
		return BadRequestResponse(nil, "empty url")
	}
	ctx := context.Background()
	response, _ := h.parser.Parse(ctx, url)

	internalInaccessible := 0
	for _, link := range response.InternalLinks {
		if !link.Accessible {
			internalInaccessible++
		}
	}
	externalInaccessible := 0
	for _, link := range response.ExternalLinks {
		if !link.Accessible {
			externalInaccessible++
		}
	}

	report := make(map[string]string)
	report["HtmlVersion"] = response.Version
	report["Title"] = response.Title
	report["H1"] = strconv.Itoa(len(response.ListH1))
	report["H2"] = strconv.Itoa(len(response.ListH2))
	report["H3"] = strconv.Itoa(len(response.ListH3))
	report["H4"] = strconv.Itoa(len(response.ListH4))
	report["H5"] = strconv.Itoa(len(response.ListH5))
	report["H6"] = strconv.Itoa(len(response.ListH6))
	report["Internal"] = strconv.Itoa(len(response.InternalLinks))
	report["InternalInaccessible"] = strconv.Itoa(internalInaccessible)
	report["External"] = strconv.Itoa(len(response.ExternalLinks))
	report["ExternalInaccessible"] = strconv.Itoa(externalInaccessible)
	report["Login"] = strconv.FormatBool(response.Login)

	tmpl := template.Must(template.ParseFiles("../static/report.html"))
	tmpl.Execute(w, model.ReportBody{Model: report})

	return nil
}
