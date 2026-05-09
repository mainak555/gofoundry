package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func AddMiddlewares(router *chi.Mux,
	allowedOrigins []string,
	middlewares ...func(http.Handler) http.Handler,
) {
	addCorsMiddleware(router, allowedOrigins)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middlewares...)

	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
}

func addCorsMiddleware(router *chi.Mux, allowedOrigins []string) {
	//allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)
}
