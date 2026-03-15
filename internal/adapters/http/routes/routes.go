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
		protected := v.Group("", opts.Protected...)

		registerSystemRoutes(public)
		switch ver {
		case "v1":
			registerPublicRoutes(public, set.V1)
			registerProtectedRoutes(protected, set.V1)
		case "v2":
			registerPublicRoutes(public, set.V2)
			registerProtectedRoutes(protected, set.V2)
		default:
			registerPublicRoutes(public, set.V1)
			registerProtectedRoutes(protected, set.V1)
		}
	}
}

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

func registerProtectedRoutes(r fiber.Router, set handlers.Set) {
	if set.User == nil {
		return
	}
	users := r.Group("/users")
	{
		users.Get("/:account_id", set.User.GetByAccountIDHandler)
		users.Post("/", nil)
		users.Put("/:id", nil)
		users.Delete("/:id", nil)
	}
}

func registerSystemRoutes(r fiber.Router) {
	r.Get("/health", func(c *fiber.Ctx) error {
		return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ok"})
	})
	r.Get("/ready", func(c *fiber.Ctx) error {
		return response.JSON(c, fiber.StatusOK, fiber.Map{"status": "ready"})
	})
}
