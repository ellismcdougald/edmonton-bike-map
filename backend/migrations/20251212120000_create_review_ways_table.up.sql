-- Create junction table for multi-way reviews
BEGIN;

-- Drop way_id from reviews since data is test-only and can be lost
ALTER TABLE reviews DROP COLUMN IF EXISTS way_id;

-- Create review_ways junction table
CREATE TABLE IF NOT EXISTS review_ways (
    review_id BIGINT NOT NULL,
    way_id BIGINT NOT NULL,
    PRIMARY KEY (review_id, way_id),
    CONSTRAINT fk_review_ways_review
        FOREIGN KEY (review_id)
        REFERENCES reviews(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_review_ways_way
        FOREIGN KEY (way_id)
        REFERENCES ways(id)
        ON DELETE CASCADE
);

-- Index to quickly look up reviews by way
CREATE INDEX IF NOT EXISTS idx_review_ways_way_id ON review_ways (way_id);

COMMIT;