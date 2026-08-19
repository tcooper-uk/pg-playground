import { useEffect, useRef, useState } from 'react'

export function usePolling<T>(
  fn: () => Promise<T>,
  intervalMs: number,
): { data: T | null; error: string | null } {
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<string | null>(null)
  const fnRef = useRef(fn)
  fnRef.current = fn

  useEffect(() => {
    let id: ReturnType<typeof setInterval>

    const tick = async () => {
      if (document.hidden) return
      try {
        setData(await fnRef.current())
        setError(null)
      } catch (e) {
        setError(e instanceof Error ? e.message : String(e))
      }
    }

    tick()
    id = setInterval(tick, intervalMs)
    return () => clearInterval(id)
  }, [intervalMs])

  return { data, error }
}
