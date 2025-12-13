-- Rollback multi-way reviews: drop junction table and restore reviews.way_id
BEGIN;

-- Drop junction table
DROP TABLE IF EXISTS review_ways;

-- Restore way_id column on reviews
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS way_id BIGINT;

-- Optionally re-add FK if desired (commented as we may not have original constraint name)
-- ALTER TABLE reviews
--     ADD CONSTRAINT fk_reviews_way
--     FOREIGN KEY (way_id)
--     REFERENCES ways(id)
--     ON DELETE SET NULL;

COMMIT;