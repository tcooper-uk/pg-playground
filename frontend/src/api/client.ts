import { Stats, WeightConfig } from './types'

export async function fetchStats(): Promise<Stats> {
  const res = await fetch('/api/stats')
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function setWeights(weights: WeightConfig): Promise<WeightConfig> {
  const res = await fetch('/api/config/weights', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(weights),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function setRate(ratePerSecond: number): Promise<void> {
  const res = await fetch('/api/config/rate', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ rate_per_second: ratePerSecond }),
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

export async function pauseReplication(): Promise<void> {
  const res = await fetch('/api/replication/pause', { method: 'POST' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}

export async function resumeReplication(): Promise<void> {
  const res = await fetch('/api/replication/resume', { method: 'POST' })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
}
