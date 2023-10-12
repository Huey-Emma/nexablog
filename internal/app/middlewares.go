package app

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"

	"nexablog/internal/models"
	"nexablog/internal/repository"
	"nexablog/internal/utils"
	"nexablog/pkg/lib"
)

func (app *App) recoverer(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		n.ServeHTTP(w, r)
	})
}

func (app *App) authenticate(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Authentication")

		pt := regexp.MustCompile(`^Bearer\s(\S+)$`)

		authheader := r.Header.Get("Authorization")

		if lib.WhiteSpace(authheader) {
			r = utils.SetUser(r, models.AnonymousUser)
			n.ServeHTTP(w, r)
			return
		}

		if !pt.Match([]byte(authheader)) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			_ = utils.WriteJson(w, http.StatusUnauthorized, lib.H[string]{
				"detail": "unauthorized",
			})
			return
		}

		token := pt.FindStringSubmatch(authheader)[1]

		user, err := app.repos.user.FindUserByToken(
			r.Context(),
			token,
			models.ScopeAuthentication,
		)

		if errors.Is(err, repository.ErrResourceNotFound) {
			w.Header().Set("WWW-Authenticate", "Bearer")
			_ = utils.WriteJson(w, http.StatusUnauthorized, lib.H[string]{
				"detail": "unauthorized",
			})
			return
		}

		if err != nil {
			_ = utils.WriteJson(w, http.StatusInternalServerError, lib.H[string]{
				"detail": "internal server error",
			})
			return
		}

		r = utils.SetUser(r, &user)
		n.ServeHTTP(w, r)
	})
}

func (app *App) requireAuth(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUser(r)

		if user.IsAnonymousUser() {
			_ = utils.WriteJson(w, http.StatusUnauthorized, lib.H[string]{
				"detail": "not authorized",
			})
			return
		}

		n.ServeHTTP(w, r)
	})
}

func (app *App) requirePermission(code string) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := utils.GetUser(r)

			permissions, err := app.repos.permission.GetUserPermission(
				r.Context(),
				user.UserID,
			)
			if err != nil {
				_ = utils.WriteJson(w, http.StatusInternalServerError, lib.H[string]{
					"detail": "internal server error",
				})
				return
			}

			if !permissions.Include(code) {
				_ = utils.WriteJson(w, http.StatusForbidden, lib.H[string]{
					"detail": "not allowed",
				})
				return
			}

			n.ServeHTTP(w, r)
		})
	}
}

func (app *App) disallowInvalidPostID(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := strconv.Atoi(chi.URLParam(r, "post-id"))

		if err != nil || postID < 1 {
			_ = utils.WriteJson(w, http.StatusNotFound, lib.H[string]{
				"detail": "post not found",
			})
			return
		}

		n.ServeHTTP(w, r)
	})
}
