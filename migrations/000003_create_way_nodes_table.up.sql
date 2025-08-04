CREATE TABLE way_nodes (
  way_id BIGINT NOT NULL,
  node_id BIGINT NOT NULL,
  sequence_id INTEGER NOT NULL,
  PRIMARY KEY (way_id, sequence_id),
  FOREIGN KEY (way_id) REFERENCES ways(id) ON DELETE CASCADE,
  FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);