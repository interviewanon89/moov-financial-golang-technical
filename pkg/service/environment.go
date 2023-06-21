package service

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/database"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"

	_ "github.com/moovfinancial/backendhiring"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger              log.Logger
	Config              *Config
	TimeService         stime.TimeService
	ZeroTrustMiddleware mux.MiddlewareFunc
	DB                  *sql.DB

	PublicRouter *mux.Router
	Shutdown     func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env == nil {
		env = &Environment{}
	}

	env.Shutdown = func() {}

	if env.Logger == nil {
		env.Logger = log.NewDefaultLogger()
	}

	if env.Config == nil {
		cfg, err := LoadConfig(env.Logger)
		if err != nil {
			return nil, err
		}

		env.Config = cfg
	}

	// db setup
	if env.DB == nil {
		db, err := database.NewAndMigrate(context.Background(), env.Logger, env.Config.Database)
		if err != nil {
			return nil, err
		}

		env.DB = db

		// Add DB closing to the Shutdown call for the Environment
		prev := env.Shutdown
		env.Shutdown = func() {
			prev()
			db.Close()
		}
	}

	if env.TimeService == nil {
		env.TimeService = stime.NewSystemTimeService()
	}

	if env.ZeroTrustMiddleware == nil {
		env.ZeroTrustMiddleware = mux.MiddlewareFunc(func(h http.Handler) http.Handler {
			return h
		})
	}

	// router
	if env.PublicRouter == nil {
		env.PublicRouter = mux.NewRouter()
	}

	env.PublicRouter.Use(env.ZeroTrustMiddleware)

	return env, nil
}

func LoadConfig(logger log.Logger) (*Config, error) {
	configService := config.NewService(logger)

	global := &GlobalConfig{}
	if err := configService.Load(global); err != nil {
		return nil, err
	}

	cfg := &global.BackendHiring

	return cfg, nil
}
