import { BrowserRouter, Routes, Route } from 'react-router-dom'
import HomePage from './pages/HomePage'
import BookingPage from './pages/BookingPage'
import CallbackPage from './pages/CallbackPage'
import AdminMoviesPage from './pages/AdminMoviesPage'

export default function App() {
  return (
    <BrowserRouter basename={import.meta.env.BASE_URL}>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/booking" element={<BookingPage />} />
        <Route path="/callback" element={<CallbackPage />} />
        <Route path="/admin/movies" element={<AdminMoviesPage />} />
      </Routes>
    </BrowserRouter>
  )
}
