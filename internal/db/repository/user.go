package repository

import (
	"fmt"
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

func (r *UserRepository) Create(chatID int64) error {
	query := `INSERT INTO users (chat_id) VALUES ($1)`
	_, err := r.db.Exec(query, chatID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
