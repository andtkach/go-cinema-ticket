import { useTimer } from '../hooks/useTimer'

export function Checkout({ session, onConfirm, onRelease, onExpire, status }) {
  const { display, urgent } = useTimer(session?.expiresAt ?? null, onExpire)

  if (!session && !status) return null

  if (!session && status) {
    return (
      <div className="checkout">
        <div className={`status-msg ${status.type}`}>{status.msg}</div>
      </div>
    )
  }

  return (
    <div className="checkout">
      <h3>Checkout</h3>
      <div className="checkout-info"><span>Seat:</span> {session.seatID}</div>
      <div className="checkout-info"><span>Movie:</span> {session.movieID}</div>
      <div className="checkout-info"><span>Session:</span> {session.sessionID.slice(0, 8)}&hellip;</div>
      <div className={`timer${urgent ? ' urgent' : ''}`}>{display}</div>
      <div className="checkout-buttons">
        <button className="btn btn--confirm" onClick={onConfirm}>Confirm</button>
        <button className="btn btn--release" onClick={onRelease}>Release</button>
      </div>
      {status && <div className={`status-msg ${status.type}`}>{status.msg}</div>}
    </div>
  )
}
