package api

import (
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func NewHandler(
	staff *service.StaffService,
	parser *service.ParserService) http.Handler {

	staffHandler := NewStaffHandler(staff)
	parserHandler := NewParserHandler(parser)
	clientHandler := NewClientHandler(parser)

	middlewareHandler := rye.NewMWHandler(rye.Config{})

	router := mux.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/health", middlewareHandler.Handle([]rye.Handler{
		staffHandler.Health,
	})).Methods(http.MethodGet)

	v1 := router.PathPrefix("/api/v1").Subrouter()

	//////////////////////////////////////////////////////////////////////////////
	// Parsing
	//////////////////////////////////////////////////////////////////////////////

	v1.Handle("/parsing/page/analyze", middlewareHandler.Handle([]rye.Handler{
		middlewareParsePayload,
		parserHandler.Parse,
	})).Methods(http.MethodPost)

	//////////////////////////////////////////////////////////////////////////////
	// Client
	//////////////////////////////////////////////////////////////////////////////

	v1.Handle("/analyzer", middlewareHandler.Handle([]rye.Handler{
		clientHandler.Form,
	})).Methods(http.MethodGet)

	v1.Handle("/analyzer", middlewareHandler.Handle([]rye.Handler{
		clientHandler.Result,
	})).Methods(http.MethodPost)

	r := cors.AllowAll()
	h := r.Handler(router)

	return h
}
