package userService

import (
 "context"
 "database/sql"
 "fmt"
 "log"
 "time"
)

// UserUpdate represents the data needed to update a user.

// PostgresUserRepository implements the UserRepository interface using PostgreSQL.
type PostgresUserRepository struct {
 db *sql.DB
}
// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
 return &PostgresUserRepository{db: db}
}
// ConnectToDB establishes a connection to the PostgreSQL database.
func ConnectToDB(dsn string) (*sql.DB, error) {
 db, err := sql.Open("postgres", dsn)
 if err != nil {
  return nil, fmt.Errorf("failed to open database: %w", err)
 }
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 if err := db.PingContext(ctx); err != nil {
  return nil, fmt.Errorf("failed to ping database: %w", err)
 }
 return db, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]User, error) {
 query := `SELECT id, name, email, created_at, updated_at FROM users`

 rows, err := r.db.QueryContext(ctx, query)
 if err != nil {
  return nil, fmt.Errorf("failed to execute query: %w", err)
 }
 defer rows.Close()

 var users []User
 for rows.Next() {
  var user User
  if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
   return nil, fmt.Errorf("failed to scan row: %w", err)
  }

  users = append(users, user)
 }

 if err := rows.Err(); err != nil {
  return nil, fmt.Errorf("error during row iteration: %w", err)
 }

 return users, nil
}

// GetByID retrieves a user from the database by ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (User, error) {
 query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

 row := r.db.QueryRowContext(ctx, query, id)

 var user User
 if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
  if err == sql.ErrNoRows {
   return User{}, fmt.Errorf("user not found with id: %d", id)
  }
  return User{}, fmt.Errorf("failed to scan row: %w", err)
 }
 return user, nil
}
// Create creates a new user in the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user UserCreate) (User, error) {
 query := `INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at`

 var createdUser User
 err := r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&createdUser.ID, &createdUser.CreatedAt, &createdUser.UpdatedAt)
 if err != nil {
  return User{}, fmt.Errorf("failed to execute query: %w", err)
 }

 createdUser.Name = user.Name
 createdUser.Email = user.Email

 return createdUser, nil
}

// Update updates a user in the database.
func (r *PostgresUserRepository) Update(ctx context.Context, id int64, user UserUpdate) (User, error) {
 query := `UPDATE users SET name = COALESCE($1, name), email = COALESCE($2, email), updated_at = NOW() WHERE id = $3 RETURNING name, email, created_at, updated_at`

 var updatedUser User
 err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, id).Scan(&updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
 if err != nil {
  if err == sql.ErrNoRows {
   return User{}, fmt.Errorf("user not found with id: %d", id)
  }
  return User{}, fmt.Errorf("failed to execute query: %w", err)
 }

    updatedUser.ID = id
 return updatedUser, nil
}
// Delete deletes a user from the database.
func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
 query := `DELETE FROM users WHERE id = $1`

 result, err := r.db.ExecContext(ctx, query, id)
 if err != nil {
  return fmt.Errorf("failed to execute query: %w", err)
 }

 rowsAffected, err := result.RowsAffected()
 if err != nil {
  log.Printf("could not get affected rows: %v", err) // Non-critical error, just log it.
 }

 if rowsAffected == 0 {
  return fmt.Errorf("user not found with id: %d", id)
 }

 return nil
}