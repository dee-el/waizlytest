package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/jackc/pgx/v5/stdlib"

	"waizlytest/config"

	commonhttpmiddleware "waizlytest/common/http/middleware"

	pgrepositories "waizlytest/repositories/pg"

	v1authhttphandler "waizlytest/services/auth/httphandlers/v1"
	jwtauthservice "waizlytest/services/auth/jwt"

	v1userhttphandler "waizlytest/services/user/httphandlers/v1"
	stduserservice "waizlytest/services/user/std"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.Read("conf.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed open conf file: %v", err))
	}

	db, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		panic(fmt.Sprintf("failed open DB: %v", err))
	}

	db.SetMaxOpenConns(cfg.DB.MaxConn)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConn)

	err = db.PingContext(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed established DB conn: %v", err))
	}

	userStorage := pgrepositories.NewUserRepository(db)

	authService, err := jwtauthservice.NewJWTAuth(userStorage, userStorage, cfg.Cert.Private, cfg.Cert.Public)
	if err != nil {
		panic(fmt.Sprintf("failed instatiate auth: %v", err))
	}

	userService := stduserservice.NewService(userStorage, userStorage)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Route("/v1", func(r chi.Router) {
		{
			hn := v1authhttphandler.NewAuthnHandler(authService)
			r.Post("/login", hn.Login())
		}

		{
			hn := v1userhttphandler.NewRegistratorHandler(userService)
			r.Post("/register", hn.Register())
		}

		// RESTy routes for "articles" resource
		r.Route("/me", func(r chi.Router) {
			{
				authMiddleware := commonhttpmiddleware.NewAuthMiddleware(authService)
				r.Use(authMiddleware.Auth)
			}

			hn := v1userhttphandler.NewMeHandler(userService)

			r.Get("/", hn.GetProfile())
			r.Put("/", hn.UpdateProfile())
		})
	})

	// start the server
	server := &http.Server{
		Addr:         cfg.Server.Port,
		ReadTimeout:  time.Second * time.Duration(cfg.Server.ReadTimeoutInSecond),
		WriteTimeout: time.Second * time.Duration(cfg.Server.WriteTimeoutInSecond),
		Handler:      r,
	}

	go func() {
		log.Printf("app server is up and running. Go to http://127.0.0.1" + cfg.Server.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	defer func() {
		db.Close()

		cancel()
	}()

	<-quit
}
