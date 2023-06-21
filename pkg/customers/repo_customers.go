package customers

import (
	"database/sql"
)

// Repository - Used for interacting identities on the data store
type CustomerRepository interface {
	Add(create Customer) (*Customer, error)
	List(tenantID string) ([]Customer, error)
	Get(tenantID string, customerID string) (*Customer, error)
	Update(update Customer) (*Customer, error)
	Delete(update Customer) (*Customer, error)
}

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func (r *customerRepo) List(tenantID string) ([]Customer, error) {
	qry := `
		SELECT 
			customers.tenant_id,
			customers.customer_id,
			customers.name,
			customers.birth_date,
			customers.email,
			customers.ssn,
			customers.created_on,
			customers.updated_on,
			customers.disabled_on
		FROM customers
		WHERE customers.tenant_id = "` + tenantID + `"
	`

	return r.queryScanCustomer(qry)
}

func (r *customerRepo) Get(tenantID string, customerID string) (*Customer, error) {
	qry := `
		SELECT 
			customers.tenant_id,
			customers.customer_id,
			customers.name,
			customers.birth_date,
			customers.email,
			customers.ssn,
			customers.created_on,
			customers.updated_on,
			customers.disabled_on
		FROM customers
		WHERE customers.tenant_id = ? 
		  AND customers.customer_id = ?
		LIMIT 1
	`

	rows, err := r.queryScanCustomer(qry, tenantID, customerID)
	if err != nil {
		return nil, err
	}

	if len(rows) != 1 {
		return nil, sql.ErrNoRows
	}

	return &rows[0], nil
}

func (r *customerRepo) Update(update Customer) (*Customer, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qry := `
		UPDATE customers
		SET
			name = ?,
			birth_date = ?,
			email = ?,
			ssn = ?,
			updated_on = ?,
			disabled_on = ?
		WHERE
			customer_id = ?
			AND tenant_id = ? 
			AND disabled_on IS NULL 
	`
	res, err := tx.Exec(qry,
		update.Name,
		update.BirthDate,
		update.Email,
		update.Ssn,
		update.UpdatedOn,
		update.DisabledOn,

		update.CustomerID,
		update.TenantID)
	if err != nil {
		return nil, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, sql.ErrNoRows
	}

	tx.Commit()

	return &update, nil
}

func (r *customerRepo) Delete(update Customer) (*Customer, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qry := `
		UPDATE customers
		SET
			updated_on = ?,
			disabled_on = ?
		WHERE
			customer_id = ? AND
			tenant_id = ? AND
			disabled_on IS NULL
	`
	res, err := tx.Exec(qry,
		update.UpdatedOn,
		update.DisabledOn,
		update.CustomerID,
		update.TenantID)
	if err != nil {
		return nil, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, sql.ErrNoRows
	}

	tx.Commit()

	return &update, nil
}

func (r *customerRepo) Add(create Customer) (*Customer, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qry := `
		INSERT INTO customers(
			tenant_id, 
			customer_id, 
			name,
			birth_date, 
			email, 
			ssn, 
			created_on, 
			updated_on, 
			disabled_on 
		) VALUES (?,?,?,?,?,?,?,?,?)
	`

	res, err := tx.Exec(qry,
		create.TenantID,
		create.CustomerID,
		create.Name,
		create.BirthDate,
		create.Email,
		create.Ssn,
		create.CreatedOn,
		create.UpdatedOn,
		create.DisabledOn,
	)
	if err != nil {
		return nil, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if cnt != 1 {
		return nil, sql.ErrNoRows
	}

	tx.Commit()

	return &create, nil
}

func (r *customerRepo) queryScanCustomer(query string, args ...interface{}) ([]Customer, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Customer{}
	for rows.Next() {
		item := Customer{}
		if err := rows.Scan(
			&item.TenantID,
			&item.CustomerID,
			&item.Name,
			&item.BirthDate,
			&item.Email,
			&item.Ssn,
			&item.CreatedOn,
			&item.UpdatedOn,
			&item.DisabledOn,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
