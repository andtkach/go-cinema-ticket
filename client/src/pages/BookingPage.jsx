import { useState, useEffect, useCallback, useRef } from 'react'
import { api } from '../api'
import { usePolling } from '../hooks/usePolling'
import { Header } from '../components/Header'
import { MovieList } from '../components/MovieList'
import { SeatGrid } from '../components/SeatGrid'
import { Checkout } from '../components/Checkout'
import { isLoggedIn, redirectToLogin, silentRefresh, getUserInfo } from '../auth'

export default function BookingPage() {
  const [ready, setReady] = useState(false)
  const [movies, setMovies] = useState([])
  const [selectedMovie, setSelectedMovie] = useState(null)
  const [seats, setSeats] = useState([])
  const [activeSession, setActiveSession] = useState(null)
  const [status, setStatus] = useState(null)
  const statusTimerRef = useRef(null)

  useEffect(() => {
    async function init() {
      if (!isLoggedIn()) {
        await redirectToLogin()
        return
      }
      silentRefresh()
      setReady(true)
    }
    init()
  }, [])

  useEffect(() => {
    if (ready) api('GET', '/movies').then(setMovies)
  }, [ready])

  const fetchSeats = useCallback(() => {
    if (!selectedMovie) return
    api('GET', `/movies/${selectedMovie.id}/seats`).then(setSeats)
  }, [selectedMovie])

  useEffect(() => { fetchSeats() }, [fetchSeats])
  usePolling(fetchSeats, 2000, !!selectedMovie)

  function showStatus(msg, type) {
    setStatus({ msg, type })
    if (statusTimerRef.current) clearTimeout(statusTimerRef.current)
    statusTimerRef.current = setTimeout(() => setStatus(null), 3000)
  }

  async function handleSelectMovie(movie) {
    if (activeSession) {
      try { await api('DELETE', `/sessions/${activeSession.sessionID}`) } catch (_) {}
      setActiveSession(null)
    }
    setSelectedMovie(movie)
    setStatus(null)
  }

  async function handleHold(seatID) {
    if (activeSession) return
    try {
      const data = await api('POST', `/movies/${selectedMovie.id}/seats/${seatID}/hold`)
      setActiveSession({
        sessionID: data.session_id,
        movieID: data.movie_id,
        seatID: data.seat_id,
        expiresAt: new Date(data.expires_at),
      })
      fetchSeats()
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  async function handleConfirm() {
    if (!activeSession) return
    try {
      await api('PUT', `/sessions/${activeSession.sessionID}/confirm`)
      setActiveSession(null)
      fetchSeats()
      showStatus('Confirmed!', 'success')
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  async function handleRelease() {
    if (!activeSession) return
    try {
      await api('DELETE', `/sessions/${activeSession.sessionID}`)
      setActiveSession(null)
      fetchSeats()
      setStatus(null)
    } catch (err) {
      showStatus(err.message, 'error')
    }
  }

  function handleExpire() {
    setActiveSession(null)
    fetchSeats()
    showStatus('Hold expired', 'error')
  }

  if (!ready) return null

  const userInfo = getUserInfo()

  return (
    <>
      <Header userInfo={userInfo} />
      <div className="booking-movies">
        <MovieList movies={movies} selectedMovie={selectedMovie} onSelect={handleSelectMovie} />
      </div>
      {selectedMovie && (
        <div className="content">
          <SeatGrid
            movie={selectedMovie}
            seats={seats}
            userID={userInfo?.sub}
            activeSession={activeSession}
            onHold={handleHold}
          />
          <Checkout
            session={activeSession}
            onConfirm={handleConfirm}
            onRelease={handleRelease}
            onExpire={handleExpire}
            status={!activeSession ? status : null}
          />
        </div>
      )}
    </>
  )
}
