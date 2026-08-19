import { useEffect, useRef } from 'react'
import { SimulatorSnapshot } from '../api/types'

interface Props {
  snapshot: SimulatorSnapshot | null
  intervalMs: number
}

interface Rate {
  rental: number
  return: number
  churn: number
  read: number
}

export function ThroughputCounters({ snapshot, intervalMs }: Props) {
  const prevRef = useRef<SimulatorSnapshot | null>(null)
  const rateRef = useRef<Rate>({ rental: 0, return: 0, churn: 0, read: 0 })

  useEffect(() => {
    if (!snapshot || !prevRef.current) {
      prevRef.current = snapshot
      return
    }
    const prev = prevRef.current
    const secs = intervalMs / 1000
    rateRef.current = {
      rental: (snapshot.rental - prev.rental) / secs,
      return: (snapshot.return - prev.return) / secs,
      churn:  (snapshot.churn  - prev.churn)  / secs,
      read:   (snapshot.read   - prev.read)   / secs,
    }
    prevRef.current = snapshot
  }, [snapshot, intervalMs])

  const rate = rateRef.current

  if (!snapshot) {
    return <div className="card"><h2>Simulator Throughput</h2><p className="empty">Simulator not running</p></div>
  }

  const ops = [
    { label: 'Rentals',  total: snapshot.rental, rate: rate.rental, colour: '#4ade80' },
    { label: 'Returns',  total: snapshot.return, rate: rate.return, colour: '#60a5fa' },
    { label: 'Churn',    total: snapshot.churn,  rate: rate.churn,  colour: '#c084fc' },
    { label: 'Reads',    total: snapshot.read,   rate: rate.read,   colour: '#facc15' },
  ]

  return (
    <div className="card">
      <h2>Simulator Throughput</h2>
      <div className="counters">
        {ops.map((op) => (
          <div key={op.label} className="counter">
            <span className="counter-label" style={{ color: op.colour }}>{op.label}</span>
            <span className="counter-total">{op.total.toLocaleString()}</span>
            <span className="counter-rate">{op.rate.toFixed(1)}/s</span>
          </div>
        ))}
      </div>
      {snapshot.errors > 0 && (
        <p className="error-note">⚠ {snapshot.errors.toLocaleString()} errors</p>
      )}
    </div>
  )
}
