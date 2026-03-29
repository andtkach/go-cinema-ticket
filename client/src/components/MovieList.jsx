export function MovieList({ movies, selectedMovie, onSelect }) {
  return (
    <div className="movies">
      {movies.map((m) => (
        <div
          key={m.id}
          className={`movie-card${selectedMovie?.id === m.id ? ' selected' : ''}`}
          onClick={() => onSelect(m)}
        >
          <h3>{m.title}</h3>
          <p>{m.rows} rows &times; {m.seats_per_row} seats</p>
        </div>
      ))}
    </div>
  )
}
