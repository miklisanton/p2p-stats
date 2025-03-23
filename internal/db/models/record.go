package models

import (
	"database/sql"
	"p2p-stats/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type Record struct {
    Id              uuid.UUID           `db:"id"`
    UserID          int64               `db:"user_id"`
    CreatedAt       time.Time           `db:"created_at"`
    DeletedAt       sql.NullTime        `db:"deleted_at"`
    Type            entities.RecordType `db:"type"`
    USDTAmount      float64             `db:"usdt_amount"`
    FiatAmount      float64             `db:"fiat_amount"`
    FiatCurrency    string              `db:"fiat_currency"`
}

type RecordUser struct {
    Id              uuid.UUID           `db:"id"`
    CreatedAt       time.Time           `db:"created_at"`
    DeletedAt       sql.NullTime        `db:"deleted_at"`
    Type            entities.RecordType `db:"type"`
    USDTAmount      float64             `db:"usdt_amount"`
    FiatAmount      float64             `db:"fiat_amount"`
    FiatCurrency    string              `db:"fiat_currency"`
    ChatID          int64               `db:"chat_id"`
    Name            string              `db:"name"`
    UserCreatedAt   time.Time           `db:"user_created_at"`
}

func ToDBRecord(record *entities.ValidatedRecord) *Record {
    return &Record{
        Id: record.Id,
        UserID: record.User.ChatID,
        CreatedAt: record.CreatedAt,
        Type: record.Type,
        USDTAmount: record.USDTAmount,
        FiatAmount: record.FiatAmount,
        FiatCurrency: record.FiatCurrency,
    }
}

func FromDBRecord(record RecordUser) *entities.Record {
    var deletedAt *time.Time
    if record.DeletedAt.Valid {
        deletedAt = &record.DeletedAt.Time
    }

    return &entities.Record{
        Id: record.Id,
        User: entities.User{
            ChatID: record.ChatID,
            CreatedAt: record.UserCreatedAt,
            Name: record.Name,
        },
        CreatedAt: record.CreatedAt,
        DeletedAt: deletedAt,
        Deleted: record.DeletedAt.Valid,
        Type: record.Type,
        USDTAmount: record.USDTAmount,
        FiatAmount: record.FiatAmount,
        FiatCurrency: record.FiatCurrency,
    }
}
