-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS cars (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    brand varchar (50) NOT NULL,
    model varchar (50) NOT NULL,
    color varchar (50) NOT NULL,
    cost numeric NOT NULL,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE cars;