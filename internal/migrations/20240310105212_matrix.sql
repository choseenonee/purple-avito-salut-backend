-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS matrix (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  microcategory_id INT,
  region_id INT,
  price INT,
  UNIQUE(name, microcategory_id, region_id)
);

CREATE TABLE IF NOT EXISTS matrix_metadata (
    id SERIAL PRIMARY KEY,
    matrix_name VARCHAR,
    timestamp TIMESTAMP,
    is_baseline bool,
    parent_matrix_name VARCHAR
);

ALTER TABLE matrix_metadata
    ADD CONSTRAINT fk_parent_matrix
        FOREIGN KEY (parent_matrix_name) REFERENCES matrix(name);

ALTER TABLE matrix_metadata
    ADD CONSTRAINT fk_matrix
        FOREIGN KEY (matrix_name) REFERENCES matrix(name) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS matrix, matrix_metadata;
-- +goose StatementEnd
