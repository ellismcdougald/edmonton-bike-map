-- Add performance indexes for nearest-way lookup optimization

-- GiST index on nodes for KNN operator (critical for fast nearest node filtering)
CREATE INDEX IF NOT EXISTS idx_nodes_point_gist ON nodes USING GIST (
  ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
);

-- Index on way_nodes (node_id) for fast candidate way collection
CREATE INDEX IF NOT EXISTS idx_way_nodes_node ON way_nodes (node_id);

-- Composite index on way_nodes (way_id, sequence_id) for fast ordered aggregation
CREATE INDEX IF NOT EXISTS idx_way_nodes_way_seq ON way_nodes (way_id, sequence_id);

ANALYZE;