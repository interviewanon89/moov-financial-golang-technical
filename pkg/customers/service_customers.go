package customers

import (
	"github.com/google/uuid"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"
)

type CustomerService interface {
	Create(tenantID string, create Customer) (*Customer, error)
	List(tenantID string) ([]Customer, error)
	Get(tenantID string, customerID string) (*Customer, error)
	Update(tenantID string, customerID string, update Customer) (*Customer, error)
	Delete(tenantID string, customerID string) error
}

func NewCustomerService(time stime.TimeService, logger log.Logger, repository CustomerRepository) (CustomerService, error) {
	return &customerService{
		time:       time,
		logger:     logger,
		repository: repository,
	}, nil
}

type customerService struct {
	time       stime.TimeService
	logger     log.Logger
	repository CustomerRepository
}

func (s *customerService) Create(tenantID string, create Customer) (*Customer, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}

	s.logger.Info().With(log.Fields{
		"Name":      log.String(create.Name),
		"BirthDate": log.StringOrNil(create.BirthDate),
		"SSN":       log.String(create.Ssn),
	}).Log("Created a new customer")

	created := Customer{
		CustomerID: uuid.New().String(),
		TenantID:   tenantID,
		CreatedOn:  s.time.Now(),
		UpdatedOn:  s.time.Now(),
		BirthDate:  create.BirthDate,
		Email:      create.Email,
		Ssn:        create.Ssn,
	}

	saved, err := s.repository.Add(created)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (s *customerService) List(tenantID string) ([]Customer, error) {
	return s.repository.List(tenantID)
}

func (s *customerService) Get(tenantID string, customerID string) (*Customer, error) {
	return s.repository.Get(tenantID, customerID)
}

func (s *customerService) Update(tenantID string, customerID string, update Customer) (*Customer, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}

	update.UpdatedOn = s.time.Now()

	_, err := s.repository.Update(update)
	if err != nil {
		return nil, err
	}

	return s.Get(tenantID, customerID)
}

func (s *customerService) Delete(tenantID string, customerID string) error {
	cur, err := s.Get(tenantID, customerID)
	if err != nil {
		return err
	}

	cur.UpdatedOn = s.time.Now()
	cur.DisabledOn = &cur.UpdatedOn

	_, err = s.repository.Delete(*cur)
	if err != nil {
		return err
	}

	return nil
}
