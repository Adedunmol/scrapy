package admin

import (
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// create admin
// get admins
// create roles
// create permissions
// assign roles
// unassign roles

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {
	adminRouter := chi.NewRouter()
	handler := Handler{
		Store: NewAdminStore(db, 5*time.Second),
	}

	// Public routes
	adminRouter.Post("/login", handler.LoginAdminHandler)
	adminRouter.Post("/", handler.CreateAdmin)

	// Protected routes
	adminRouter.Group(func(protected chi.Router) {
		protected.Use(helpers.AuthMiddleware)

		//// Admins
		protected.Get("/", handler.GetAdmins)

		//// Roles
		protected.Post("/roles", handler.CreateRole)
		protected.Post("/roles/batch", handler.BatchCreateRoles)
		protected.Get("/roles", handler.GetRoles)

		// Permissions
		protected.Post("/permissions", handler.CreatePermission)
		protected.Post("/permissions/batch", handler.BatchCreatePermissions)
		protected.Get("/roles/{id}/permissions", handler.GetPermissions)
	})

	r.Mount("/admins", adminRouter)
}
