package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/admin"
	"github.com/moov-io/base/log"

	_ "github.com/moovfinancial/backendhiring"
)

// RunServers - Boots up all the servers and awaits till they are stopped.
func (env *Environment) RunServers(terminationListener chan error) func() {

	adminServer := bootAdminServer(terminationListener, env.Logger, env.Config.Servers.Admin)

	_, shutdownPublicServer := bootHTTPServer("public", env.PublicRouter, terminationListener, env.Logger, env.Config.Servers.Public)

	return func() {
		adminServer.Shutdown()
		shutdownPublicServer()
	}
}

func bootHTTPServer(name string, routes *mux.Router, errs chan<- error, logger log.Logger, config HTTPConfig) (*http.Server, func()) {

	loggedHandler := RequestLogger(logger, routes, "http")

	// Create main HTTP server
	serve := &http.Server{
		Addr:    config.Bind.Address,
		Handler: loggedHandler,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:       false,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start main HTTP server
	go func() {
		logger.Info().Log(fmt.Sprintf("%s listening on %s", name, config.Bind.Address))
		if err := serve.ListenAndServe(); err != nil {
			errs <- logger.Fatal().LogErrorf("problem starting http: %w", err).Err()
		}
	}()

	shutdownServer := func() {
		if err := serve.Shutdown(context.Background()); err != nil {
			logger.Fatal().LogErrorf("shutting down: %v", err)
		}
	}

	return serve, shutdownServer
}

func bootAdminServer(errs chan<- error, logger log.Logger, config HTTPConfig) *admin.Server {
	adminServer := admin.NewServer(config.Bind.Address)

	go func() {
		logger.Info().Log(fmt.Sprintf("listening on %s", adminServer.BindAddr()))
		if err := adminServer.Listen(); err != nil {
			errs <- logger.Fatal().LogErrorf("problem starting admin http: %w", err).Err()
		}
	}()

	return adminServer
}

func RequestLogger(logger log.Logger, inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := log.Fields{
			"request_method": log.String(r.Method),
			"require_uri":    log.String(r.RequestURI),
			"route_name":     log.String(name),
			"response_time":  log.TimeDuration(time.Since(start)),
		}

		defer func() {
			if r := recover(); r != nil {
				ctx["panic"] = log.String(fmt.Sprintf("%v", r))
			}

			logger.With(ctx).Info().Logf("%s %s %s", r.Method, r.RequestURI, name)
		}()

		inner.ServeHTTP(w, r)
	})
}
