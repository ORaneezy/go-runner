CREATE TABLE pipeline_steps (
    id SERIAL PRIMARY KEY,
    name TEXT,
    sequence_order INTEGER,
    command TEXT,
    pipeline_id INTEGER REFERENCES pipelines(id)
);