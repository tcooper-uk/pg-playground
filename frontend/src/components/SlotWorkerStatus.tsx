import { SlotHealth, WorkerInfo } from '../api/types'

interface Props {
  slots: SlotHealth[] | null
  workers: WorkerInfo | null
}

export function SlotWorkerStatus({ slots, workers }: Props) {
  return (
    <div className="card">
      <h2>Slots &amp; Workers</h2>

      <h3>Replication Slots</h3>
      {!slots || slots.length === 0 ? (
        <p className="empty">No logical slots</p>
      ) : (
        slots.map((s) => (
          <div key={s.slot_name} className="info-row">
            <span className="dot" style={{ background: s.active ? '#4ade80' : '#ef4444' }} />
            <strong>{s.slot_name}</strong>
            <span className="label">{s.active ? 'active' : 'inactive'}</span>
            {s.retained_bytes != null && (
              <span className="label">{s.retained_bytes.toLocaleString()} bytes retained</span>
            )}
          </div>
        ))
      )}

      <h3>WAL Senders</h3>
      {!workers?.senders || workers.senders.length === 0 ? (
        <p className="empty">No senders</p>
      ) : (
        workers.senders.map((s, i) => (
          <div key={i} className="info-row">
            <span className="dot" style={{ background: s.state === 'streaming' ? '#4ade80' : '#facc15' }} />
            <strong>{s.application_name}</strong>
            <span className="label">{s.state}</span>
            <span className="label mono">sent {s.sent_lsn}</span>
            <span className="label mono">replay {s.replay_lsn}</span>
          </div>
        ))
      )}

      <h3>Receivers</h3>
      {!workers?.receivers || workers.receivers.length === 0 ? (
        <p className="empty">No receivers</p>
      ) : (
        workers.receivers.map((r, i) => (
          <div key={i} className="info-row">
            <span className="dot" style={{ background: '#4ade80' }} />
            <strong>{r.sub_name}</strong>
            <span className="label mono">received {r.received_lsn}</span>
            <span className="label mono">end {r.latest_end_lsn}</span>
          </div>
        ))
      )}
    </div>
  )
}
