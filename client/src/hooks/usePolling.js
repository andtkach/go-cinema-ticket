import { useEffect } from 'react'

export function usePolling(fn, intervalMs, active) {
  useEffect(() => {
    if (!active) return
    const id = setInterval(fn, intervalMs)
    return () => clearInterval(id)
  }, [fn, active, intervalMs])
}
