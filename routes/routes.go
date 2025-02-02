package routes

import (
	"database/sql"
	"io/fs"

	"github.com/DeanWard/erugo/handlers"
	"github.com/DeanWard/erugo/middleware"
	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all routes for the application
func RegisterRoutes(router *mux.Router, database *sql.DB, embeddedFS fs.FS) {
	// Auth routes
	registerAuthRoutes(router, database)

	// Share routes
	registerShareRoutes(router, database)

	// User routes
	registerUserRoutes(router, database)

	// Health check
	registerHealthRoutes(router)

	// Frontend routes
	registerFrontendRoutes(router, embeddedFS)
}

func registerAuthRoutes(router *mux.Router, database *sql.DB) {
	//POST /api/auth/login - login a user
	router.Handle("/api/auth/login",
		handlers.LoginHandler(database),
	).Methods("POST")

	//POST /api/auth/refresh - refresh a user's token
	router.Handle("/api/auth/refresh",
		handlers.RefreshTokenHandler(database),
	).Methods("POST")

	//POST /api/auth/logout - logout a user
	router.Handle("/api/auth/logout",
		handlers.LogoutHandler(database),
	).Methods("POST")
}

func registerShareRoutes(router *mux.Router, database *sql.DB) {
	//POST /api/shares - create a new share
	router.Handle("/api/shares",
		middleware.JwtMiddleware(
			handlers.CreateShareHandler(database),
		),
	).Methods("POST")

	//GET /api/shares/{longId} - get a share by longId
	router.Handle("/api/shares/{longId}",
		handlers.GetShareHandler(database),
	).Methods("GET")

	//GET /api/shares/{longId}/download - download a share by longId
	router.Handle("/api/shares/{longId}/download",
		handlers.DownloadShareHandler(database),
	).Methods("GET")
}

func registerUserRoutes(router *mux.Router, database *sql.DB) {
	//GET /api/users - get all users
	router.Handle("/api/users",
		middleware.JwtMiddleware(
			middleware.AdminMiddleware(
				handlers.GetUsersHandler(database),
			),
		),
	).Methods("GET")

	//POST /api/users - create a new user
	router.Handle("/api/users",
		middleware.JwtMiddleware(
			middleware.AdminMiddleware(
				handlers.CreateUserHandler(database),
			),
		),
	).Methods("POST")

}

func registerHealthRoutes(router *mux.Router) {
	router.Handle("/api/health",
		handlers.HealthCheckHandler(),
	).Methods("GET")
}

func registerFrontendRoutes(router *mux.Router, embeddedFS fs.FS) {
	router.PathPrefix("/").Handler(handlers.ServeFrontendHandler(embeddedFS))
}
