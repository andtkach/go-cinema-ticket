import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { api } from '../api'
import { isLoggedIn, isAdmin, getUserInfo } from '../auth'
import { Header } from '../components/Header'
import { MovieForm } from '../components/MovieForm'

export default function AdminMoviesPage() {
  const navigate = useNavigate()
  const [ready, setReady] = useState(false)
  const [movies, setMovies] = useState([])
  const [creating, setCreating] = useState(false)
  const [editing, setEditing] = useState(null)
  const [status, setStatus] = useState(null)
  const statusTimerRef = useRef(null)

  useEffect(() => {
    if (!isLoggedIn() || !isAdmin()) {
      navigate('/', { replace: true })
      return
    }
    setReady(true)
  }, [navigate])

  useEffect(() => {
    if (ready) loadMovies()
  }, [ready])

  async function loadMovies() {
    try {
      const data = await api('GET', '/movies')
      setMovies(data || [])
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  function showStatus(msg, type) {
    setStatus({ msg, type })
    if (statusTimerRef.current) clearTimeout(statusTimerRef.current)
    statusTimerRef.current = setTimeout(() => setStatus(null), 3000)
  }

  async function handleCreate(data) {
    try {
      await api('POST', '/movies', data)
      setCreating(false)
      await loadMovies()
      showStatus('Movie created', 'success')
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  async function handleUpdate(data) {
    try {
      await api('PUT', `/movies/${editing.id}`, data)
      setEditing(null)
      await loadMovies()
      showStatus('Movie updated', 'success')
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  async function handleDelete(id) {
    try {
      await api('DELETE', `/movies/${id}`)
      await loadMovies()
      showStatus('Movie deleted', 'success')
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  if (!ready) return null

  const userInfo = getUserInfo()

  return (
    <>
      <Header userInfo={userInfo} />
      <div className="admin-toolbar">
        <h2 className="home-section-title" style={{ marginBottom: 0 }}>Manage Movies</h2>
        {!creating && (
          <button className="btn btn--confirm admin-add-btn" onClick={() => { setCreating(true); setEditing(null) }}>
            + Add Movie
          </button>
        )}
      </div>

      {creating && (
        <div className="movie-card" style={{ margin: '0 auto 1.5rem', maxWidth: '360px' }}>
          <h3 style={{ marginBottom: '0.4rem', color: 'var(--accent)' }}>New Movie</h3>
          <MovieForm
            onSubmit={handleCreate}
            onCancel={() => setCreating(false)}
            submitLabel="Create"
          />
        </div>
      )}

      {status && (
        <p className={`status-msg ${status.type}`} style={{ marginBottom: '1rem' }}>{status.msg}</p>
      )}

      <div className="movies">
        {movies.map(m => (
          <div key={m.id} className="movie-card movie-card--admin">
            {editing?.id === m.id ? (
              <>
                <h3 style={{ marginBottom: '0.4rem', color: 'var(--accent)' }}>Edit Movie</h3>
                <MovieForm
                  initial={m}
                  onSubmit={handleUpdate}
                  onCancel={() => setEditing(null)}
                  submitLabel="Update"
                />
              </>
            ) : (
              <>
                <h3>{m.title}</h3>
                <p>{m.rows} rows &times; {m.seats} seats</p>
                <div className="movie-card-actions">
                  <button
                    className="btn btn--confirm"
                    onClick={() => { setEditing(m); setCreating(false) }}
                  >
                    Edit
                  </button>
                  <button
                    className="btn btn--release"
                    onClick={() => handleDelete(m.id)}
                  >
                    Delete
                  </button>
                </div>
              </>
            )}
          </div>
        ))}
      </div>

      {movies.length === 0 && !creating && (
        <p className="empty-state">No movies yet. Add one above.</p>
      )}
    </>
  )
}
