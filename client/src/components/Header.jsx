import { logout } from '../auth'

export function Header({ userInfo }) {
  const display = userInfo?.name || userInfo?.email || userInfo?.sub || 'unknown'
  return (
    <header>
      <h1>Cinema Booking</h1>
      <div className="user-id">
        <span>user: {display}</span>
        <button onClick={logout}>Logout</button>
      </div>
    </header>
  )
}
