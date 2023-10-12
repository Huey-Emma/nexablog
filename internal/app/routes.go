package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"nexablog/internal/handlers"
	"nexablog/internal/services/permission"
	"nexablog/internal/services/post"
	"nexablog/internal/services/token"
	"nexablog/internal/services/user"
	"nexablog/internal/utils"
	"nexablog/pkg/lib"
)

func (app *App) loadRoutes() {
	api := chi.NewRouter()

	api.Use(
		app.recoverer,
		middleware.Logger,
		app.authenticate,
	)

	api.Get("/", app.wrap(handlers.Welcome))

	api.Route("/users", app.loadUserRoutes)
	api.Route("/tokens", app.loadTokenRoutes)
	api.Route("/posts", app.loadPostRoutes)

	app.mux.Mount("/api", api)
}

func (app *App) loadPostRoutes(r chi.Router) {
	postSvc := post.NewService(app.repos.post)

	h := handlers.Post{
		PostSvc: postSvc,
	}

	r.Get("/", app.wrap(h.FindAllPosts))

	r.With(
		app.requireAuth,
		app.requirePermission("posts:write"),
	).Post("/", app.wrap(h.CreatePost))

	r.With(
		app.disallowInvalidPostID,
	).Get("/{post-id:^[0-9]+}", app.wrap(h.FindPostByID))

	r.With(
		app.requireAuth,
		app.disallowInvalidPostID,
	).Delete("/{post-id:^[0-9]+}", app.wrap(h.DeletePostByID))

	r.With(
		app.requireAuth,
		app.requirePermission("posts:write"),
		app.disallowInvalidPostID,
	).Put("/{post-id:[0-9]+}", app.wrap(h.UpdatePostByID))
}

func (app *App) loadTokenRoutes(r chi.Router) {
	userSvc := user.NewService(app.repos.user)
	tokenSvc := token.NewService(app.repos.token)

	h := handlers.Token{
		UserSvc:  userSvc,
		TokenSvc: tokenSvc,
	}

	r.Post("/authenticate", app.wrap(h.GetToken))
}

func (app *App) loadUserRoutes(r chi.Router) {
	userSvc := user.NewService(app.repos.user)
	permissionSvc := permission.NewService(app.repos.permission)
	postSvc := post.NewService(app.repos.post)

	h := handlers.User{
		UserSvc:       userSvc,
		PermissionSvc: permissionSvc,
		PostSvc:       postSvc,
	}

	r.Post("/", app.wrap(h.Register))

	r.With(
		app.requireAuth,
	).Get("/me", app.wrap(h.GetMe))
}

func (app *App) wrap(f utils.ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		var apiError *utils.ApiError

		switch {
		case errors.As(err, &apiError):
			_ = utils.WriteJson(w, apiError.Code(), lib.H[string]{
				"detail": apiError.Error(),
			})
			return
		case err != nil:
			log.Println(err)
			_ = utils.WriteJson(w, http.StatusInternalServerError, lib.H[string]{
				"detail": "internal server error",
			})
		}
	}
}
