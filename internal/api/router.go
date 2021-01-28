package api

import (
	"net/http"
	"os"
	"path"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/log"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/InVisionApp/rye"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewHandler(
	staff *service.StaffService,
	parser *service.ParserService) http.Handler {

	staffHandler := NewStaffHandler(staff)
	parserHandler := NewParserHandler(parser)
	clientHandler := NewClientHandler(nil)

	middlewareHandler := rye.NewMWHandler(rye.Config{})

	router := mux.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/health", middlewareHandler.Handle([]rye.Handler{
		staffHandler.Health,
	})).Methods(http.MethodGet)

	workDir, _ := os.Getwd()
	log.Infof("working_dir: %v", workDir)
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir(path.Join(workDir, "./swagger/")))))

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

	v1.Handle("/parsing/page/analyze", middlewareHandler.Handle([]rye.Handler{
		clientHandler.Client,
	})).Methods(http.MethodGet)

	r := cors.AllowAll()
	h := r.Handler(router)

	return h
}
