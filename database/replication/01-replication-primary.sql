CREATE USER replicator WITH REPLICATION LOGIN PASSWORD 'replicator';
GRANT CONNECT ON DATABASE dvdrental TO replicator;
GRANT pg_read_all_data TO replicator;
CREATE PUBLICATION dvdrental_pub FOR ALL TABLES;
