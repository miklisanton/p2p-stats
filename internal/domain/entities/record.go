package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RecordType string

const (
    Buy     RecordType = "BUY"
    Sell    RecordType = "SELL"
)

type Record struct {
    Id              uuid.UUID
    CreatedAt       time.Time
    DeletedAt       *time.Time
    Deleted         bool
    Type            RecordType
    USDTAmount      float64
    FiatAmount      float64
    FiatCurrency    string
    UserID          int64
}

func NewRecord(t RecordType, usdtAmount, fiatAmount float64, fiatCurrency string, userID int64) *Record {
    return &Record{
        Id: uuid.New(),
        CreatedAt: time.Now(),
        Type: t,
        USDTAmount: usdtAmount,
        FiatAmount: fiatAmount,
        FiatCurrency: fiatCurrency,
        UserID: userID,
    }
}

func (r *Record) validate() error {
    if r.UserID == 0 {
        return errors.New("invalid user")
    }
    if r.USDTAmount <= 0 {
        return errors.New("invalid USDT amount")
    }
    if r.FiatAmount <= 0 {
        return errors.New("invalid fiat amount")
    }
    if len(r.FiatCurrency) != 3 {
        return errors.New("invalid fiat currency")
    }
    return nil
}

func (r *Record) String() string {
    return fmt.Sprintf("%s %.0fUSDT for %.0f%s", r.Type, r.USDTAmount, r.FiatAmount, r.FiatCurrency)
}
