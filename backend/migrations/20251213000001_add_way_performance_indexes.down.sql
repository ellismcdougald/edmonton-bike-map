-- Drop performance indexes for nearest-way lookup

DROP INDEX IF EXISTS idx_nodes_point_gist;
DROP INDEX IF EXISTS idx_way_nodes_node;
DROP INDEX IF EXISTS idx_way_nodes_way_seq;
