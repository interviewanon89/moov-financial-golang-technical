package customers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/log"
)

type CustomerController interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewCustomerController(logger log.Logger, service CustomerService) CustomerController {
	return &customerController{
		logger:  logger,
		service: service,
	}
}

type customerController struct {
	logger  log.Logger
	service CustomerService
}

func (c customerController) AppendRoutes(router *mux.Router) *mux.Router {
	router.
		Name("Customer.create").
		Methods("PUT").
		Path("/customers").
		HandlerFunc(c.create)

	router.
		Name("Customer.list").
		Methods("GET").
		Path("/customers").
		HandlerFunc(c.list)

	router.
		Name("Customer.get").
		Methods("GET").
		Path("/customers/{ID}").
		HandlerFunc(c.get)

	router.
		Name("Customer.update").
		Methods("PUT").
		Path("/customer/{ID}").
		HandlerFunc(c.update)

	router.
		Name("Customer.delete").
		Methods("DELETE").
		Path("/customers/{ID}").
		HandlerFunc(c.delete)

	return router
}

func (c *customerController) GetTenantID(r *http.Request) (string, error) {
	tenID := r.Header.Get("X-Tenant-ID")
	if tenID == "" {
		return "", errors.New("Missing tenantID")
	}
	return tenID, nil
}

func (c *customerController) create(w http.ResponseWriter, r *http.Request) {
	tenantID, err := c.GetTenantID(r)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	create := Customer{}
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	result, err := c.service.Create(tenantID, create)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	jsonResponse(w, result)
}

func (c *customerController) list(w http.ResponseWriter, r *http.Request) {
	tenantID, err := c.GetTenantID(r)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	result, err := c.service.List(tenantID)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	jsonResponse(w, result)
}

func (c *customerController) get(w http.ResponseWriter, r *http.Request) {
	tenantID, err := c.GetTenantID(r)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	params := mux.Vars(r)
	customerID := params["ID"]

	result, err := c.service.Get(customerID, tenantID)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	jsonResponse(w, result)
}

func (c *customerController) update(w http.ResponseWriter, r *http.Request) {
	tenantID, err := c.GetTenantID(r)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	params := mux.Vars(r)
	customerID := params["ID"]

	update := Customer{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	result, err := c.service.Update(tenantID, customerID, update)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	jsonResponse(w, result)
}

func (c *customerController) delete(w http.ResponseWriter, r *http.Request) {
	tenantID, err := c.GetTenantID(r)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	params := mux.Vars(r)
	customerID := params["ID"]

	err = c.service.Delete(tenantID, customerID)
	if err != nil {
		errorResponse(w, err, c.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
