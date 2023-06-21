CREATE TABLE customers (
    tenant_id           VARCHAR(36) NOT NULL,
    customer_id         VARCHAR(36) NOT NULL,

    name                VARCHAR(255) NOT NULL,
    birth_date          VARCHAR(10),
    email               VARCHAR(255) NOT NULL,
    ssn                 VARCHAR(11) NOT NULL,

    created_on          TIMESTAMP NOT NULL,
    updated_on          TIMESTAMP NOT NULL,
    disabled_on         TIMESTAMP,

    CONSTRAINT customer_pk PRIMARY KEY (tenant_id, customer_id)
);