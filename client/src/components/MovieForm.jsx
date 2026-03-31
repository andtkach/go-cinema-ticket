import { useState } from 'react'

export function MovieForm({ initial = {}, onSubmit, onCancel, submitLabel = 'Save' }) {
  const [title, setTitle] = useState(initial.title ?? '')
  const [rows, setRows] = useState(initial.rows ?? 5)
  const [seats, setSeats] = useState(initial.seats ?? 8)

  function handleSubmit(e) {
    e.preventDefault()
    onSubmit({ title, rows: Number(rows), seats: Number(seats) })
  }

  return (
    <form className="movie-form" onSubmit={handleSubmit}>
      <label>
        Title
        <input
          type="text"
          value={title}
          onChange={e => setTitle(e.target.value)}
          required
        />
      </label>
      <label>
        Rows
        <input
          type="number"
          min="1"
          value={rows}
          onChange={e => setRows(e.target.value)}
          required
        />
      </label>
      <label>
        Seats per row
        <input
          type="number"
          min="1"
          value={seats}
          onChange={e => setSeats(e.target.value)}
          required
        />
      </label>
      <div className="checkout-buttons">
        <button type="submit" className="btn btn--confirm">{submitLabel}</button>
        <button type="button" className="btn btn--release" onClick={onCancel}>Cancel</button>
      </div>
    </form>
  )
}
