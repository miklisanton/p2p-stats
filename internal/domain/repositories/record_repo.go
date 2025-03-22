package repositories

import "p2p-stats/internal/domain/entities"

type RecordRepo interface {
    Create(record *entities.Record) (*entities.Record, error)
    GetByUserID(userID int64) ([]entities.Record, error)
}

