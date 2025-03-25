package repository

import (
	"fmt"
	"p2p-stats/internal/db/models"
	"p2p-stats/internal/domain/entities"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type RecordRepository struct {
	db *sqlx.DB
}

func NewRecordRepository(db *sqlx.DB) *RecordRepository {
	return &RecordRepository{
		db: db,
	}
}


func (r *RecordRepository) Create(record *entities.ValidatedRecord) error {
    query := `
    INSERT INTO records (id, user_id, created_at, type, usdt_amount, fiat_amount, fiat_currency)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`

    _, err := r.db.Exec(query, record.Id, record.UserID, record.CreatedAt, record.Type, record.USDTAmount, record.FiatAmount, record.FiatCurrency)
    if err != nil {
        return fmt.Errorf("failed to create record: %w", err)
    }
    return nil
}

func (r *RecordRepository) GetByUserID(userID int64, from, to *time.Time) ([]entities.Record, error) {
    query := `
    SELECT r.* FROM records r
    WHERE r.user_id = $1`
    if from != nil {
        query += " AND r.created_at >= $2"
    }
    if to != nil {
        query += " AND r.created_at <= $3"
        t := to.Add(24 * time.Hour)
        to = &t
        log.Info().Msgf("to: %v", to)
    }

    var records []models.Record
    err := r.db.Select(&records, query, userID, from, to)
    if err != nil {
        return nil, fmt.Errorf("failed to get records: %w", err)
    }
    
    var result []entities.Record
    for _, record := range records {
        result = append(result, *models.FromDBRecord(record))
    }
    return result, nil
}

func (r *RecordRepository) Delete(id string) error {
    query := `DELETE FROM records WHERE id = $1`
    _, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete record: %w", err)
    }
    return nil
}
