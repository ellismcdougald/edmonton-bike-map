CREATE TABLE reviews (
  id SERIAL PRIMARY KEY,
  way_id BIGINT NOT NULL REFERENCES ways(id) ON DELETE CASCADE,
  user_id BIGINT, -- Users not implemented yet. Will become a foreign key later
  rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
  comment TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);