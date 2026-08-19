select * from pg_catalog.pg_user pu;
select * from pg_catalog.pg_publication pp;
select * from pg_catalog.pg_replication_slots prs;
select * from pg_catalog.pg_subscription ps;

CREATE USER replicator WITH REPLICATION LOGIN PASSWORD 'replicator';
GRANT CONNECT ON DATABASE dvdrental TO replicator;
GRANT pg_read_all_data TO replicator;
CREATE PUBLICATION dvdrental_pub FOR ALL TABLES;

-- WAL sender state and lag
SELECT
    application_name,
    state,
    pg_wal_lsn_diff(pg_current_wal_lsn(), sent_lsn)   AS unsent_bytes,
    pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS lag_bytes,
    write_lag,
    flush_lag,
    replay_lag
FROM pg_stat_replication;

-- Replication slot — confirms slot is active and shows retained WAL
SELECT
    slot_name,
    active,
    pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn) AS retained_bytes
FROM pg_replication_slots;


insert into actor (first_name, last_name , last_update ) values ('Tom', 'Cooper', clock_timestamp());