export function Header({ userID }) {
  return (
    <header>
      <h1>Cinema Booking</h1>
      <div className="user-id">user: {userID}</div>
    </header>
  )
}
