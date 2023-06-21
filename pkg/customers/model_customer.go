package customers

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Customer struct {
	// UUID v4
	TenantID string `json:"tenantID,omitempty"`
	// UUID v4
	CustomerID string  `json:"customerID,omitempty"`
	Name       string  `json:"name,omitempty"`
	BirthDate  *string `json:"birthDate,omitempty"`
	// Email Address
	Email      string     `json:"email,omitempty"`
	Ssn        string     `json:"ssn,omitempty"`
	CreatedOn  time.Time  `json:"createdOn,omitempty"`
	UpdatedOn  time.Time  `json:"updatedOn,omitempty"`
	DisabledOn *time.Time `json:"disabledOn,omitempty"`
}

func (a Customer) Validate() error {
	// Ozzo validation: https://github.com/go-ozzo/ozzo-validation#validating-a-simple-value
	return validation.ValidateStruct(&a,
		validation.Field(&a.CustomerID, validation.Required, is.UUID),
		validation.Field(&a.TenantID, validation.Required, is.UUID),
	)
}
