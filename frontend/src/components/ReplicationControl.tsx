import { useState } from 'react'
import { SubscriptionStatus } from '../api/types'
import { pauseReplication, resumeReplication } from '../api/client'

interface Props {
  subscriptions: SubscriptionStatus[] | null
}

export function ReplicationControl({ subscriptions }: Props) {
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const enabled = subscriptions?.some((s) => s.enabled) ?? true

  const handle = async (action: () => Promise<void>) => {
    setBusy(true)
    setError(null)
    try {
      await action()
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="card">
      <h2>Replication Control</h2>

      <div className="repl-status">
        <span className="dot" style={{ background: enabled ? '#4ade80' : '#ef4444' }} />
        <span className="repl-status-text" style={{ color: enabled ? '#4ade80' : '#ef4444' }}>
          {enabled ? 'Active' : 'Paused'}
        </span>
        {subscriptions?.map((s) => (
          <span key={s.name} className="label mono">{s.name}</span>
        ))}
      </div>

      <p className="repl-hint">
        Pause the subscription to let WAL build up on the primary, then resume to watch the lag drain on the chart.
      </p>

      <div className="repl-buttons">
        <button
          className="repl-btn repl-btn--pause"
          onClick={() => handle(pauseReplication)}
          disabled={busy || !enabled}
        >
          Pause
        </button>
        <button
          className="repl-btn repl-btn--resume"
          onClick={() => handle(resumeReplication)}
          disabled={busy || enabled}
        >
          Resume
        </button>
      </div>

      {error && <p className="error-note">{error}</p>}
    </div>
  )
}
