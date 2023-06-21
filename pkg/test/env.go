package test

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/backendhiring/pkg/service"
)

type TestEnvironment struct {
	Assert     *require.Assertions
	StaticTime stime.StaticTimeService
	TenantID   string

	service.Environment
}

func NewEnvironment(t *testing.T, router *mux.Router) *TestEnvironment {
	assert := require.New(t)
	logger := log.NewNopLogger() //log.NewDefaultLogger()

	testEnv := &TestEnvironment{}
	testEnv.Assert = assert
	testEnv.StaticTime = stime.NewStaticTimeService()

	testEnv.TenantID = uuid.New().String()
	mw := mux.MiddlewareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Tenant-ID") == "" {
				r.Header.Set("X-Tenant-ID", testEnv.TenantID)
			}

			h.ServeHTTP(rw, r)
		})
	})

	env, err := service.NewEnvironment(&service.Environment{
		Logger:              logger,
		TimeService:         testEnv.StaticTime,
		ZeroTrustMiddleware: mw,
		PublicRouter:        router,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(env.Shutdown)

	testEnv.Environment = *env
	return testEnv
}

func SQLiteDBPath(t *testing.T) string {
	dbPath, err := ioutil.TempFile("", "sqlite-test.*.db")
	if err != nil {
		t.Fatalf("sqlite temp file failure: %v", err)
	}

	// Cleanup the database after the test has ran
	t.Cleanup(func() {
		os.Remove(dbPath.Name())
	})

	return dbPath.Name()
}
