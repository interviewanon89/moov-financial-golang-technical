package customers_test

import (
	"database/sql"
	"testing"

	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/moov-io/base/database"
	"github.com/moovfinancial/backendhiring/pkg/customers"
	"github.com/stretchr/testify/require"
)

func NewCustomer() customers.Customer {
	m := NewTestCustomer(nil)
	return m
}

func Test_Customer_AddAndGet(t *testing.T) {
	CustomerTestEachDatabase(t, func(t *testing.T, repository customers.CustomerRepository) {
		a := require.New(t)

		model := NewCustomer()

		added, err := repository.Add(model)
		a.Nil(err)
		a.Equal(model, *added)

		found, err := repository.Get(added.TenantID, added.CustomerID)
		a.Nil(err)

		a.Nil(added.DisabledOn)

		a.Equal(*added, *found)

		badTenantID := uuid.New().String()
		_, err = repository.Get(badTenantID, added.CustomerID)
		if err != sql.ErrNoRows {
			t.Fatal(err)
		}
	})
}

func Test_Customer_List(t *testing.T) {
	CustomerTestEachDatabase(t, func(t *testing.T, repository customers.CustomerRepository) {
		a := require.New(t)

		added, err := repository.Add(NewCustomer())
		a.Nil(err)

		tenantID := added.TenantID

		// Add noise and other invites on other tenants
		_, _ = repository.Add(NewCustomer())
		_, _ = repository.Add(NewCustomer())
		_, _ = repository.Add(NewCustomer())
		_, _ = repository.Add(NewCustomer())

		found, err := repository.List(tenantID)
		a.Nil(err)
		a.Len(found, 1)
		a.Equal(*added, found[0])

		badTenantID := uuid.New().String()
		found, err = repository.List(badTenantID)
		a.Nil(err)
		a.Empty(found)
	})
}

func Test_Customer_Update(t *testing.T) {
	CustomerTestEachDatabase(t, func(t *testing.T, repository customers.CustomerRepository) {
		a := require.New(t)

		added, err := repository.Add(NewCustomer())
		a.Nil(err)

		tenantID := added.TenantID
		updated := *added
		updated.UpdatedOn = time.Now().UTC()

		// @TODO add some valid changes here

		saved, err := repository.Update(updated)
		a.Nil(err)
		a.Equal(updated, *saved)

		found, err := repository.Get(tenantID, updated.CustomerID)
		a.Nil(err)
		a.Equal(updated, *found)

		badUpdate := updated
		badUpdate.TenantID = uuid.New().String()
		_, err = repository.Update(badUpdate)
		a.Equal(err, sql.ErrNoRows)
	})
}

func Test_Customer_Delete(t *testing.T) {
	CustomerTestEachDatabase(t, func(t *testing.T, repository customers.CustomerRepository) {
		a := require.New(t)

		added, err := repository.Add(NewCustomer())
		a.Nil(err)

		tenantID := added.TenantID
		updated := *added
		updated.UpdatedOn = time.Now().UTC()
		updated.DisabledOn = &updated.UpdatedOn

		// @TODO add some valid changes here

		saved, err := repository.Delete(updated)
		a.Nil(err)
		a.Equal(updated, *saved)

		// We can retrieve deleted items by specifically asking for them.
		got, err := repository.Get(tenantID, updated.CustomerID)
		a.Nil(err)
		a.Equal(updated, *got)

		// Don't list anything thats been deleted
		listed, err := repository.List(tenantID)
		a.Nil(err)
		a.Empty(listed)

		badUpdate := updated
		badUpdate.TenantID = uuid.New().String()
		_, err = repository.Update(badUpdate)
		a.Equal(err, sql.ErrNoRows)
	})
}

func CustomerTestEachDatabase(t *testing.T, run func(t *testing.T, repository customers.CustomerRepository)) {
	cases := map[string]*sql.DB{
		"sqlite": database.CreateTestSQLiteDB(t).DB,
	}

	for k, db := range cases {
		t.Run(k, func(t *testing.T) {
			repo := customers.NewCustomerRepository(db)
			run(t, repo)
		})
	}
}
