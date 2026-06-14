CREATE TYPE pipeline_status AS ENUM ('waiting', 'running', 'success', 'failure');

CREATE TABLE pipeline_runs (
    id SERIAL PRIMARY KEY,
    status pipeline_status,
    pipeline_id INTEGER REFERENCES pipelines(id)
);