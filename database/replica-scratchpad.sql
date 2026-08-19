select * from pg_catalog.pg_user pu;
select * from pg_catalog.pg_publication pp;
select * from pg_catalog.pg_replication_slots prs;
select * from pg_catalog.pg_subscription ps;


CREATE SUBSCRIPTION dvdrental_sub
	CONNECTION 'host=dvdrental-primary port=5432 dbname=dvdrental user=replicator password=replicator' 
	PUBLICATION dvdrental_pub;


-- Subscription state — shows last received LSN and worker status
SELECT
    subname,
    received_lsn,
    latest_end_lsn,
    latest_end_time
FROM pg_stat_subscription;

-- per-table sync sate
SELECT
    srrelid::regclass AS table_name,
    CASE srsubstate
        WHEN 'i' THEN 'initializing'
        WHEN 'd' THEN 'copying data'
        WHEN 'f' THEN 'finished copy'
        WHEN 's' THEN 'synchronized'
        WHEN 'r' THEN 'ready / streaming'
    END AS state
FROM pg_subscription_rel;

-- copy status
SELECT
    relid::regclass AS table_name,
    command,
    tuples_processed,
    tuples_excluded
FROM pg_stat_progress_copy;