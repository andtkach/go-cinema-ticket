import { useState, useEffect } from 'react'

export function useTimer(expiresAt, onExpire) {
  const [remaining, setRemaining] = useState(0)

  useEffect(() => {
    if (!expiresAt) { setRemaining(0); return }
    const tick = () => {
      const r = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))
      setRemaining(r)
      if (r <= 0) onExpire()
    }
    tick()
    const id = setInterval(tick, 1000)
    return () => clearInterval(id)
  }, [expiresAt]) // eslint-disable-line react-hooks/exhaustive-deps

  const mins = String(Math.floor(remaining / 60)).padStart(2, '0')
  const secs = String(remaining % 60).padStart(2, '0')
  return { display: `${mins}:${secs}`, urgent: remaining < 60 }
}
