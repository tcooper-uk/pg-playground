import { useRef, useState } from 'react'

export interface LagSample {
  t: number
  lagBytes: number
  replayLagMs: number
}

const MAX_SAMPLES = 60

export function useLagHistory() {
  const [history, setHistory] = useState<LagSample[]>([])
  const bufRef = useRef<LagSample[]>([])

  const push = (lagBytes: number, replayLagMs: number) => {
    const sample: LagSample = { t: Date.now(), lagBytes, replayLagMs }
    const next = [...bufRef.current, sample].slice(-MAX_SAMPLES)
    bufRef.current = next
    setHistory(next)
  }

  return { history, push }
}
