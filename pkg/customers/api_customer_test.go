package customers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/moovfinancial/backendhiring/pkg/customers"
	"github.com/stretchr/testify/require"
)

func Test_Customer_CreateAPI(t *testing.T) {
	s := CustomerTestSetup(t)
	m := NewTestCustomer(s.Env.TimeService)

	found, resp, err := clientCustomerCreate(s, m)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)

	s.Assert.NotEmpty(found.CustomerID)
	s.Assert.NotEmpty(found.CreatedOn)
	s.Assert.NotEmpty(found.TenantID)

	s.Assert.Equal(m.Name, found.Name)
	s.Assert.Equal(m.Email, found.Email)
	s.Assert.Equal(m.BirthDate, found.BirthDate)
	s.Assert.Equal(m.Ssn, found.Ssn)
}

func Test_Customer_Validation_Email(t *testing.T) {
	customer := NewTestCustomer(nil)

	customer.Email = ""
	err := customer.Validate()
	require.Error(t, err)
}

func Test_Customer_ListAPI(t *testing.T) {
	s := CustomerTestSetup(t)

	// Lets add a random one in here
	m := addFuzzedCustomer(s)

	found, resp, err := clientCustomerList(s)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)
	s.Assert.Len(found, 1)

	s.Assert.Contains(found, m)
}

func Test_Customer_ListAPI_NotFound(t *testing.T) {
	s := CustomerTestSetup(t)

	found, resp, err := clientCustomerList(s)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)
	s.Assert.Len(found, 0)
}

func Test_Customer_DeleteAPI(t *testing.T) {
	s := CustomerTestSetup(t)

	// Lets add a random one in here
	m := addFuzzedCustomer(s)

	// Call delete to disable the customer.
	resp, err := clientCustomerDelete(s, m.CustomerID)
	s.Assert.Nil(err)
	s.Assert.Equal(204, resp.StatusCode)

	// Direct fetch of the model should return so people can review it.
	disabled, resp, err := clientCustomerGet(s, m.CustomerID)
	s.Assert.NotNil(resp)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)

	s.Assert.Equal(s.Env.StaticTime.Now(), disabled.UpdatedOn)

	// Listing of the model resource should not show a disabled model.
	found, resp, err := clientCustomerList(s)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)
	s.Assert.Len(found, 0)
}

func Test_Customer_DeleteAPI_NotFound(t *testing.T) {
	s := CustomerTestSetup(t)

	resp, _ := clientCustomerDelete(s, uuid.New().String())
	s.Assert.NotNil(resp)
	s.Assert.Equal(404, resp.StatusCode)
}

func Test_Customer_GetAPI(t *testing.T) {
	s := CustomerTestSetup(t)

	// Lets add a random one in here
	m := addFuzzedCustomer(s)

	found, resp, err := clientCustomerGet(s, m.CustomerID)

	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)
	s.Assert.Equal(m, found)
}

func Test_Customer_GetAPI_NotFound(t *testing.T) {
	s := CustomerTestSetup(t)

	_, resp, _ := clientCustomerGet(s, uuid.New().String())
	s.Assert.NotNil(resp)
	s.Assert.Equal(404, resp.StatusCode)
}

func Test_Customer_UpdateAPI(t *testing.T) {
	s := CustomerTestSetup(t)

	// Lets add a random one in here to the repository
	m := addFuzzedCustomer(s)

	// fast forward time
	s.Env.StaticTime.Add(time.Hour)

	// Fuzz a new customer api object to fuzz in new values
	bd := "1980/03/31"
	now := s.Env.StaticTime.Now()
	updates := customers.Customer{
		TenantID:   uuid.NewString(),
		CustomerID: uuid.NewString(),
		Name:       "Jane Doe",
		Email:      "jane.doe@moov.io",
		BirthDate:  &bd,
		Ssn:        "111-222-3333",
		CreatedOn:  now.Add(time.Hour),
		UpdatedOn:  now.Add(time.Hour),
		DisabledOn: &now,
	}

	// Do the update
	updated, resp, err := clientCustomerUpdate(s, m.CustomerID, updates)
	s.Assert.Nil(err)
	s.Assert.NotNil(resp)
	s.Assert.Equal(200, resp.StatusCode)

	// These shouldn't change.
	s.Assert.Equal(m.CustomerID, updated.CustomerID)
	s.Assert.Equal(s.Env.TenantID, updated.TenantID)
	s.Assert.Equal(m.CreatedOn, updated.CreatedOn)
	s.Assert.Equal(s.Env.TimeService.Now(), updated.UpdatedOn)
	s.Assert.Nil(updated.DisabledOn)

	// Change based on session and when request was received.
	s.Assert.Equal(s.Env.StaticTime.Now().Unix(), updated.UpdatedOn.Unix())

	// These change
	s.Assert.Equal(updates.Name, updated.Name)
	s.Assert.Equal(updates.Email, updated.Email)
	s.Assert.Equal(updates.BirthDate, updated.BirthDate)
	s.Assert.Equal(updates.Ssn, updated.Ssn)

	// Lets fetch it fresh and check that it matches
	found, resp, err := clientCustomerGet(s, m.CustomerID)
	s.Assert.NotNil(resp)
	s.Assert.Nil(err)
	s.Assert.Equal(200, resp.StatusCode)
	s.Assert.Equal(updated, found)
}

func Test_Customer_UpdateAPI_NotFound(t *testing.T) {
	s := CustomerTestSetup(t)

	updates := NewTestCustomer(s.Env.TimeService)

	// If we attempt to update a customer that doesn't exist it should 404
	doesNotExistID := uuid.New().String()
	_, resp, _ := clientCustomerUpdate(s, doesNotExistID, updates)
	s.Assert.NotNil(resp)
	s.Assert.Equal(404, resp.StatusCode)
}

// Generate a random customer and insert it into the database and return it.
func addFuzzedCustomer(s CustomerTestScope) customers.Customer {
	m := NewTestCustomer(s.Env.TimeService)
	m.TenantID = s.Env.TenantID
	m.CreatedOn = s.Env.TimeService.Now()
	m.UpdatedOn = s.Env.TimeService.Now()
	_, err := s.Repository.Add(m)
	s.Assert.Nil(err)

	return customers.Customer{
		TenantID:   m.TenantID,
		CustomerID: m.CustomerID,
		Name:       m.Name,
		BirthDate:  m.BirthDate,
		Email:      m.Email,
		Ssn:        m.Ssn,
		CreatedOn:  m.CreatedOn,
		UpdatedOn:  m.UpdatedOn,
		DisabledOn: m.DisabledOn,
	}
}

// These function calls below are generated from the OpenAPI specifications in the ./api/yml

func clientCustomerCreate(s CustomerTestScope, create customers.Customer) (customers.Customer, *http.Response, error) {
	cus := customers.Customer{}
	res := s.MakeCall(s.MakeRequest("POST", "/customers", &create), &cus)
	return cus, res, nil
}

func clientCustomerList(s CustomerTestScope) ([]customers.Customer, *http.Response, error) {
	cus := []customers.Customer{}
	res := s.MakeCall(httptest.NewRequest("GET", "/customers", nil), &cus)
	return cus, res, nil
}

func clientCustomerGet(s CustomerTestScope, customerID string) (customers.Customer, *http.Response, error) {
	cus := customers.Customer{}
	res := s.MakeCall(httptest.NewRequest("GET", "/customers/"+customerID, nil), &cus)
	return cus, res, nil
}

func clientCustomerUpdate(s CustomerTestScope, customerID string, updates customers.Customer) (customers.Customer, *http.Response, error) {
	cus := customers.Customer{}
	res := s.MakeCall(s.MakeRequest("PUT", "/customers/"+customerID, &updates), &cus)
	return cus, res, nil
}

func clientCustomerDelete(s CustomerTestScope, customerID string) (*http.Response, error) {
	res := s.MakeCall(s.MakeRequest("DELETE", "/customers/"+customerID, nil), nil)
	return res, nil
}
