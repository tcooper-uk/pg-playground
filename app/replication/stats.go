package replication

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Lag struct {
	ApplicationName string         `json:"application_name"`
	State           string         `json:"state"`
	LagBytes        *int64         `json:"lag_bytes"`
	WriteLag        *time.Duration `json:"write_lag_ns"`
	FlushLag        *time.Duration `json:"flush_lag_ns"`
	ReplayLag       *time.Duration `json:"replay_lag_ns"`
}

type TableState struct {
	TableName string `json:"table_name"`
	State     string `json:"state"`
}

type SlotHealth struct {
	SlotName      string `json:"slot_name"`
	Active        bool   `json:"active"`
	RetainedBytes *int64 `json:"retained_bytes"`
}

type WorkerInfo struct {
	Senders   []SenderInfo   `json:"senders"`
	Receivers []ReceiverInfo `json:"receivers"`
}

type SenderInfo struct {
	ApplicationName string `json:"application_name"`
	State           string `json:"state"`
	SentLSN         string `json:"sent_lsn"`
	ReplayLSN       string `json:"replay_lsn"`
}

type ReceiverInfo struct {
	SubName      string    `json:"sub_name"`
	ReceivedLSN  string    `json:"received_lsn"`
	LatestEndLSN string    `json:"latest_end_lsn"`
	LatestEndTime time.Time `json:"latest_end_time"`
}

func GetLag(ctx context.Context, primary *pgxpool.Pool) ([]Lag, error) {
	rows, err := primary.Query(ctx, `
		SELECT
			application_name,
			state,
			pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) AS lag_bytes,
			extract(epoch FROM write_lag)::bigint,
			extract(epoch FROM flush_lag)::bigint,
			extract(epoch FROM replay_lag)::bigint
		FROM pg_stat_replication`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Lag
	for rows.Next() {
		var l Lag
		var writeLagNs, flushLagNs, replayLagNs *int64
		if err := rows.Scan(&l.ApplicationName, &l.State, &l.LagBytes, &writeLagNs, &flushLagNs, &replayLagNs); err != nil {
			return nil, err
		}
		if writeLagNs != nil {
			d := time.Duration(*writeLagNs) * time.Second
			l.WriteLag = &d
		}
		if flushLagNs != nil {
			d := time.Duration(*flushLagNs) * time.Second
			l.FlushLag = &d
		}
		if replayLagNs != nil {
			d := time.Duration(*replayLagNs) * time.Second
			l.ReplayLag = &d
		}
		results = append(results, l)
	}
	return results, rows.Err()
}

func GetSubscriptionState(ctx context.Context, replica *pgxpool.Pool) ([]TableState, error) {
	rows, err := replica.Query(ctx, `
		SELECT
			srrelid::regclass::text AS table_name,
			CASE srsubstate
				WHEN 'i' THEN 'initializing'
				WHEN 'd' THEN 'copying data'
				WHEN 'f' THEN 'finished copy'
				WHEN 's' THEN 'synchronized'
				WHEN 'r' THEN 'ready'
			END AS state
		FROM pg_subscription_rel
		ORDER BY table_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TableState
	for rows.Next() {
		var t TableState
		if err := rows.Scan(&t.TableName, &t.State); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, rows.Err()
}

func GetSlotHealth(ctx context.Context, primary *pgxpool.Pool) ([]SlotHealth, error) {
	rows, err := primary.Query(ctx, `
		SELECT
			slot_name,
			active,
			pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn) AS retained_bytes
		FROM pg_replication_slots
		WHERE slot_type = 'logical'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SlotHealth
	for rows.Next() {
		var s SlotHealth
		if err := rows.Scan(&s.SlotName, &s.Active, &s.RetainedBytes); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, rows.Err()
}

func GetWorkerInfo(ctx context.Context, primary, replica *pgxpool.Pool) (*WorkerInfo, error) {
	info := &WorkerInfo{}

	senderRows, err := primary.Query(ctx, `
		SELECT application_name, state, sent_lsn::text, replay_lsn::text
		FROM pg_stat_replication`)
	if err != nil {
		return nil, err
	}
	defer senderRows.Close()
	for senderRows.Next() {
		var s SenderInfo
		if err := senderRows.Scan(&s.ApplicationName, &s.State, &s.SentLSN, &s.ReplayLSN); err != nil {
			return nil, err
		}
		info.Senders = append(info.Senders, s)
	}
	if err := senderRows.Err(); err != nil {
		return nil, err
	}

	receiverRows, err := replica.Query(ctx, `
		SELECT subname, received_lsn::text, latest_end_lsn::text, latest_end_time
		FROM pg_stat_subscription
		WHERE received_lsn IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer receiverRows.Close()
	for receiverRows.Next() {
		var r ReceiverInfo
		if err := receiverRows.Scan(&r.SubName, &r.ReceivedLSN, &r.LatestEndLSN, &r.LatestEndTime); err != nil {
			return nil, err
		}
		info.Receivers = append(info.Receivers, r)
	}
	if err := receiverRows.Err(); err != nil {
		return nil, err
	}

	return info, nil
}
