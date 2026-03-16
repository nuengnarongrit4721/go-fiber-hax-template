package routes

import (
	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
)

type Options struct {
	Versions  []string
	Public    []fiber.Handler
	Protected []fiber.Handler
}

func Register(app *fiber.App, set handlers.VersionedSet, opts Options) {
	api := app.Group("/api")

	versions := opts.Versions
	if len(versions) == 0 {
		versions = []string{"v1"}
	}

	for _, ver := range versions {
		v := api.Group("/" + ver)
		public := v.Group("", opts.Public...)

		registerSystemRoutes(public)
		switch ver {
		case "v1":
			registerPublicRoutes(public, set.V1)
			registerProtectedRoutes(v, set.V1, opts.Protected...)
		case "v2":
			registerPublicRoutes(public, set.V2)
			registerProtectedRoutes(v, set.V2, opts.Protected...)
		default:
			registerPublicRoutes(public, set.V1)
			registerProtectedRoutes(v, set.V1, opts.Protected...)
		}
	}
}

/* NO Middleware */
func registerPublicRoutes(r fiber.Router, set handlers.Set) {
	if set.Auth == nil {
		return
	}
	auth := r.Group("/auth")
	{
		auth.Post("/login", set.Auth.Login)
		auth.Post("/register", set.Auth.Register)
	}
}

/* Middleware */
func registerProtectedRoutes(r fiber.Router, set handlers.Set, middleware ...fiber.Handler) {
	if set.User == nil {
		return
	}
	users := r.Group("/users", middleware...)
	{
		users.Get("/:account_id", set.User.GetByAccountIDHandler)
	}
}

/* System Routes */
func registerSystemRoutes(r fiber.Router) {
	r.Get("/health", func(c *fiber.Ctx) error {
		return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ok"})
	})
	r.Get("/ready", func(c *fiber.Ctx) error {
		return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ready"})
	})
}
