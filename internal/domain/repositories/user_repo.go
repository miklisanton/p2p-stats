package repositories

import "p2p-stats/internal/domain/entities"

type UserRepository interface {
    Create(user *entities.ValidatedUser) (*entities.User, error)
    GetAll() ([]entities.User, error)
}
