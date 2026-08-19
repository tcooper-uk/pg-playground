# PG-PLAYGROUND

A collection of test experiments backed by postgres.

# Prereqs

- Docker
- psql
- Go
- Node.js 18+
- Internet
- Patience
- An interest in postgres

# Experiments

- Logical replication between instances
- Schema migrations with pgroll

# Foundations

A postgres sample database and some applications to simulate a real world workloads. Samples are based on the collection from [Neon](https://github.com/neondatabase/postgres-sample-dbs).

# Setup

## Start the containers

```bash
make up
```

Brings up `dvdrental-primary` on port `5432` and `dvdrental-replica` on port `5433`. The replica waits for the primary to be healthy before starting.

## Seed the primary

```bash
make seed
```

Loads `database/dvdrental.sql` into the primary.

## Set up logical replication

Run the three steps in order:

```bash
make replica-init-1   # primary: create replication user and publication
make replica-init-2   # replica: apply schema
make replica-init-3   # replica: create subscription
```

Or run all three at once:

```bash
make replica-init
```

After step 3 the replica connects to the primary and begins copying data. Monitor progress on the replica:

```sql
SELECT srrelid::regclass AS table_name,
       CASE srstate
           WHEN 'i' THEN 'initializing'
           WHEN 'd' THEN 'copying data'
           WHEN 'f' THEN 'finished copy'
           WHEN 's' THEN 'synchronized'
           WHEN 'r' THEN 'ready / streaming'
       END AS state
FROM pg_subscription_rel;
```

Check replication lag on the primary:

```sql
SELECT application_name,
       state,
       pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS lag_bytes,
       replay_lag
FROM pg_stat_replication;
```

## Tear everything down

```bash
make reset
```

Stops containers and wipes both data directories for a clean slate.

# Go App

The app connects to both the primary and replica, runs a configurable workload simulator, and exposes replication metrics over HTTP.

## Build and run

```bash
make app-build
make app-run
```

The app reads `app/config.yaml`. Key settings:

```yaml
simulator:
  enabled: true
  rate_per_second: 10   # total ops/sec across all workers
  workers: 4
  weights:              # relative frequency of each op type
    rental: 40
    return: 30
    customer_churn: 10
    read: 20
```

## API endpoints

All endpoints return JSON and are served on port `8080`.

| Endpoint | Description |
|---|---|
| `GET /health` | Liveness check |
| `GET /stats` | All metrics combined |
| `GET /stats/lag` | Replication lag bytes + time from primary |
| `GET /stats/tables` | Per-table sync state from replica |
| `GET /stats/slots` | Replication slot health from primary |
| `GET /stats/workers` | WAL sender + subscription worker info |
| `GET /stats/simulator` | Simulator op counters |

# Dashboard

A React + Vite frontend that visualises replication metrics with live updates (polling every second).

## Install and run

```bash
make frontend-install   # first time only
make frontend-dev       # starts Vite dev server on http://localhost:5173
```

The dev server proxies `/api/*` to the Go app on `:8080`, so both must be running simultaneously.

## What's on screen

- **Lag chart** — time-series of replication lag bytes and replay lag, updated each second
- **Throughput counters** — per-op totals and ops/sec from the simulator
- **Slots & workers** — slot active status, retained WAL bytes, and sender/receiver LSNs
- **Table sync grid** — all tables colour-coded by sync state (useful to watch during initial copy)