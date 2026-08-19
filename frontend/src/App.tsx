import { useEffect } from 'react'
import { fetchStats } from './api/client'
import { usePolling } from './hooks/usePolling'
import { useLagHistory } from './hooks/useLagHistory'
import { LagChart } from './components/LagChart'
import { TableSyncGrid } from './components/TableSyncGrid'
import { SlotWorkerStatus } from './components/SlotWorkerStatus'
import { ThroughputCounters } from './components/ThroughputCounters'
import './styles.css'

const POLL_MS = 1000

export default function App() {
  const { data, error } = usePolling(fetchStats, POLL_MS)
  const { history, push } = useLagHistory()

  useEffect(() => {
    if (!data?.lag?.length) return
    const row = data.lag[0]
    push(
      row.lag_bytes ?? 0,
      row.replay_lag_ns != null ? row.replay_lag_ns / 1_000_000 : 0,
    )
  }, [data])

  return (
    <div className="app">
      <header>
        <h1>PG Playground — Replication Dashboard</h1>
        {error && <span className="error-badge">API error: {error}</span>}
      </header>
      <div className="grid">
        <div className="span-2">
          <LagChart history={history} />
        </div>
        <ThroughputCounters snapshot={data?.simulator ?? null} intervalMs={POLL_MS} />
        <SlotWorkerStatus slots={data?.slots ?? null} workers={data?.workers ?? null} />
        <div className="span-2">
          <TableSyncGrid tables={data?.tables ?? null} />
        </div>
      </div>
    </div>
  )
}
