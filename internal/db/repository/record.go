package repository

import (
	"fmt"
	"p2p-stats/internal/db/models"
	"p2p-stats/internal/domain/entities"

	"github.com/jmoiron/sqlx"
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

    _, err := r.db.Exec(query, record.Id, record.User.ChatID, record.CreatedAt, record.Type, record.USDTAmount, record.FiatAmount, record.FiatCurrency)
    if err != nil {
        return fmt.Errorf("failed to create record: %w", err)
    }
    return nil
}

func (r *RecordRepository) GetByUserID(userID int64) ([]entities.Record, error) {
    query := `
    SELECT r.*, u.chat_id, u.name, u.created_at as user_created_at FROM records r
    JOIN users u ON r.user_id = u.chat_id
    WHERE r.user_id = $1`

    var records []models.RecordUser
    err := r.db.Select(&records, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get records: %w", err)
    }
    
    var result []entities.Record
    for _, record := range records {
        result = append(result, *models.FromDBRecord(record))
    }
    return result, nil
}
