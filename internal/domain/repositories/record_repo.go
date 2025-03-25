package repositories

import (
	"p2p-stats/internal/domain/entities"
	"time"
)

type RecordRepo interface {
    Create(record *entities.ValidatedRecord) error
    GetByUserID(userID int64, from, to *time.Time) ([]entities.Record, error)
    Delete(recordID string) error
}

