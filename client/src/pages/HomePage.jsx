import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Header } from '../components/Header'
import { fetchPublic } from '../api'

export default function HomePage() {
  const [cinemas, setCinemas] = useState([])

  useEffect(() => {
    fetchPublic('/cinemas').then(data => setCinemas(data || []))
  }, [])

  return (
    <>
      <Header />
      <div className="home-hero">
        <h2>Welcome to Cinema Booking</h2>
        <p>Pick a cinema, choose your seat, and confirm in seconds.</p>
        <Link to="/booking" className="btn btn--confirm home-cta">Book a Seat</Link>
      </div>
      {cinemas.length > 0 && (
        <>
          <h2 className="home-section-title">Our Cinemas</h2>
          <div className="movies">
            {cinemas.map(c => (
              <div key={c.id} className="movie-card">
                <h3>{c.name}</h3>
                <p>{c.location}</p>
                <p>{c.screens} screen{c.screens !== 1 ? 's' : ''}</p>
                <Link to="/booking" className="movie-card-link">Book now →</Link>
              </div>
            ))}
          </div>
        </>
      )}
    </>
  )
}
