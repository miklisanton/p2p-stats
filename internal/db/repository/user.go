package repository

import (
	"fmt"
	"p2p-stats/internal/db/models"
	"p2p-stats/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *entities.ValidatedUser) error {
    u := models.ToDBUser(user)
	query := `INSERT INTO users (chat_id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, u.ChatID, u.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetAll() ([]entities.User, error) {
    query := `SELECT * FROM users`
    var users []models.User
    err := r.db.Select(&users, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }

    var result []entities.User
    for _, user := range users {
        result = append(result, *models.FromDBUser(&user))
    }
    return result, nil
}
