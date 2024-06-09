package larago

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (l *Larago) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	if l.Debug {
		router.Use(middleware.Logger)
	}
	router.Use(middleware.Recoverer)

	return router
}
