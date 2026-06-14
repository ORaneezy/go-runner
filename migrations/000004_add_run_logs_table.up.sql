CREATE TABLE run_logs (
  id SERIAL PRIMARY KEY,
  step_id INTEGER REFERENCES pipeline_steps(id),
  message TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  run_id INTEGER REFERENCES pipeline_runs(id)
);