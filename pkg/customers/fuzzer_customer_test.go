package customers_test

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/moov-io/base/stime"
	"github.com/moovfinancial/backendhiring/pkg/customers"
)

func NewTestCustomer(times stime.TimeService) customers.Customer {
	if times == nil {
		times = stime.NewStaticTimeService()
	}

	ssn := fmt.Sprintf("%d-%d-%d", ((rand.Int() % 899) + 100), ((rand.Int() % 89) + 10), ((rand.Int() % 8999) + 1000))

	bd := times.Now().Format("2006/02/01")
	return customers.Customer{
		TenantID:   uuid.NewString(),
		CustomerID: uuid.NewString(),
		Name:       "Joe J Doe",
		Ssn:        ssn,
		Email:      "john.doe@moov.io",
		BirthDate:  &bd,
	}
}
