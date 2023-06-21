package customers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	fuzz "github.com/google/gofuzz"
	"github.com/moovfinancial/backendhiring/pkg/customers"
	"github.com/moovfinancial/backendhiring/pkg/test"
	"github.com/stretchr/testify/require"
)

type CustomerTestScope struct {
	T      *testing.T
	Assert *require.Assertions
	Fuzzer *fuzz.Fuzzer
	Env    *test.TestEnvironment

	Repository customers.CustomerRepository
	Service    customers.CustomerService

	Router *mux.Router
}

func CustomerTestSetup(t *testing.T) CustomerTestScope {
	a := require.New(t)

	router := mux.NewRouter()
	testEnv := test.NewEnvironment(t, router)

	// These can be replaced with whats in the `testEnv` created above.
	repository := customers.NewCustomerRepository(testEnv.DB)
	service, _ := customers.NewCustomerService(testEnv.TimeService, testEnv.Logger, repository)
	controller := customers.NewCustomerController(testEnv.Logger, service)

	controller.AppendRoutes(router)

	return CustomerTestScope{
		T:          t,
		Assert:     a,
		Env:        testEnv,
		Repository: repository,
		Service:    service,
		Router:     router,
	}
}

func (s CustomerTestScope) MakeRequest(method string, target string, body interface{}) *http.Request {
	jsonBody := bytes.Buffer{}
	if body != nil {
		json.NewEncoder(&jsonBody).Encode(body)
	}

	return httptest.NewRequest(method, target, &jsonBody)
}

func (s CustomerTestScope) MakeCall(req *http.Request, body interface{}) *http.Response {
	rec := httptest.NewRecorder()
	s.Router.ServeHTTP(rec, req)
	res := rec.Result()

	if body != nil {
		json.NewDecoder(res.Body).Decode(&body)
	}

	return res
}
