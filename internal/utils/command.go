package utils

import (
	"errors"
	"p2p-stats/internal/domain/entities"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ParseList(args string) (*time.Time, *time.Time, error) {
    parts := strings.Fields(args)
    if len(parts) > 2 {
        return nil, nil, errors.New("invalid list, expected 0, 1 or 2 arguments")
    }

    if len(parts) == 0 {
        return nil, nil, nil
    }

    var from, to *time.Time
    if len(parts) > 0 {
        f, err := time.Parse("02-01-2006", parts[0])
        from = &f
        if err != nil {
            return nil, nil, errors.New("invalid from date, expected format DD-MM-YYYY")
        }
    }

    if len(parts) == 2 {
        t, err := time.Parse("02-01-2006", parts[1])
        to = &t
        if err != nil {
            return nil, nil, errors.New("invalid to date, expected format DD-MM-YYYY")
        }
    }

    return from, to, nil
}

func ParseRecord(t entities.RecordType, args string) (*entities.Record, error) {
    // parse record from string
    parts := strings.Fields(args)
    if len(parts) != 3 {
        return nil, errors.New("invalid record, expected 3 arguments")
    }
    usdtAmount, err := strconv.ParseFloat(parts[0], 64)
    if err != nil {
        return nil, errors.New("invalid USDT amount, expected number, e.g. 0.5")
    }
    fiatAmount, err := strconv.ParseFloat(parts[1], 64)
    if err != nil {
        return nil, errors.New("invalid fiat amount, expected number, e.g. 500")
    }
    fiatCurrency := strings.ToUpper(parts[2])
    if len(fiatCurrency) != 3 {
        return nil, errors.New("invalid fiat currency, expected 3 letters, e.g. USD")
    }
    return &entities.Record{
        Id: uuid.New(),
        CreatedAt: time.Now(),
        Type: t,
        USDTAmount: usdtAmount,
        FiatAmount: fiatAmount,
        FiatCurrency: fiatCurrency,
    }, nil
}
