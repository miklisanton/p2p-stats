package repositories

import "p2p-stats/internal/domain/entities"

type RecordRepo interface {
    Create(record *entities.ValidatedRecord) (*entities.Record, error)
    GetByUserID(userID int64) ([]entities.Record, error)
}

