import { Link } from 'react-router-dom'
import { logout, isAdmin } from '../auth'

export function Header({ userInfo }) {
  const display = userInfo?.name || userInfo?.email || userInfo?.sub || 'unknown'
  return (
    <header>
      <h1><Link to="/" className="header-home-link">Cinema Booking</Link></h1>
      {userInfo ? (
        <div className="header-user">
          {isAdmin() && (
            <Link to="/admin/movies" className="btn btn--confirm header-signin-btn">Admin</Link>
          )}
          <div className="user-id">
            <span>user: {display}</span>
            <button onClick={logout}>Logout</button>
          </div>
        </div>
      ) : (
        <Link to="/booking" className="btn btn--confirm header-signin-btn">Sign In</Link>
      )}
    </header>
  )
}
