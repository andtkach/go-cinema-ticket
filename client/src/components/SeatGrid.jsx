const ROW_LABELS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'

export function SeatGrid({ movie, seats, userID, activeSession, onHold }) {
  const statusMap = Object.fromEntries((seats || []).map((s) => [s.seat_id, s]))

  return (
    <div className="screen-area">
      <div className="screen-label">Screen</div>
      <div className="screen-bar" />
      <div className="seat-grid">
        {Array.from({ length: movie.rows }, (_, r) => (
          <div key={r} className="seat-row">
            <div className="row-label">{ROW_LABELS[r]}</div>
            {Array.from({ length: movie.seats }, (_, s) => {
              const seatID = ROW_LABELS[r] + (s + 1)
              const info = statusMap[seatID]
              let cls = 'seat'
              if (info?.confirmed) cls += ' seat--confirmed'
              else if (info?.booked && info.user_id === userID) cls += ' seat--held-mine'
              else if (info?.booked) cls += ' seat--held-other'
              const disabled = !!activeSession || (info?.booked || info?.confirmed)
              return (
                <button
                  key={seatID}
                  className={cls}
                  disabled={disabled}
                  onClick={() => !activeSession && !info?.booked && !info?.confirmed && onHold(seatID)}
                >
                  {s + 1}
                </button>
              )
            })}
            <div className="row-label">{ROW_LABELS[r]}</div>
          </div>
        ))}
      </div>
      <div className="legend">
        {[
          { label: 'Available', color: 'var(--available)' },
          { label: 'Your hold', color: 'var(--held-mine)' },
          { label: 'Other hold', color: 'var(--held-other)' },
          { label: 'Confirmed', color: 'var(--confirmed)' },
        ].map(({ label, color }) => (
          <div key={label} className="legend-item">
            <div className="legend-swatch" style={{ background: color }} />
            {label}
          </div>
        ))}
      </div>
    </div>
  )
}
