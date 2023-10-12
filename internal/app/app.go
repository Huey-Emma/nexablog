package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"nexablog/config"
	"nexablog/db"
	"nexablog/internal/repository/permission"
	"nexablog/internal/repository/post"
	"nexablog/internal/repository/token"
	"nexablog/internal/repository/user"
)

type App struct {
	cfg      *config.Config
	database *db.DB
	mux      *chi.Mux
	repos    *repos
}

type repos struct {
	user       user.Repo
	token      token.Repo
	post       post.Repo
	permission permission.Repo
}

func New(cfg *config.Config, database *db.DB) *App {
	app := &App{
		cfg:      cfg,
		database: database,
		mux:      chi.NewRouter(),
	}

	app.loadRepos()
	app.loadRoutes()

	return app
}

func (app *App) loadRepos() {
	r := &repos{
		user:       user.NewRepo(app.database),
		token:      token.NewRepo(app.database),
		post:       post.NewRepo(app.database),
		permission: permission.NewRepo(app.database),
	}

	app.repos = r
}

func (app *App) StartAndRun(ctx context.Context) error {
	errch := make(chan error)

	if err := app.database.PingDB(ctx); err != nil {
		return fmt.Errorf("could not ping db: %w", err)
	}

	defer func() {
		log.Println("database is closing")
		err := app.database.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	server := &http.Server{
		Addr:         "0.0.0.0:" + app.cfg.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Handler:      app.mux,
	}

	log.Println("app is starting")

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			errch <- err
		}
		errch <- nil
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Println("app is shutting down")

		err := server.Shutdown(ctx)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}
