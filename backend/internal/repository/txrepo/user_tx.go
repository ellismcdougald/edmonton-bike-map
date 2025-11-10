package txrepo

import (
    "database/sql"

    "github.com/ellismcdougald/edmonton-bike-map/internal/models"
    "github.com/ellismcdougald/edmonton-bike-map/internal/repository"
)

// TxUserRepository provides user operations within an existing transaction.
type TxUserRepository struct {
    Tx *sql.Tx
}

// NewTxUserRepository returns a UserRepository that operates using the provided transaction.
func NewTxUserRepository(tx *sql.Tx) repository.UserRepository {
    return &TxUserRepository{Tx: tx}
}

// GetByUsername retrieves a user by username within the transaction.
func (r *TxUserRepository) GetByUsername(username string) (*models.User, error) {
    var u models.User
    err := r.Tx.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
        Scan(&u.ID, &u.Username, &u.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &u, nil
}

// Create inserts a new user within the transaction.
func (r *TxUserRepository) Create(user *models.User) error {
    _, err := r.Tx.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
    return err
}

// UsernameExists checks if a username exists within the transaction.
func (r *TxUserRepository) UsernameExists(username string) (bool, error) {
    var exists bool
    err := r.Tx.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
    return exists, err
}
