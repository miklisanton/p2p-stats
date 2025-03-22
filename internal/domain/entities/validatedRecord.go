package entities

type ValidatedRecord struct {
    Record
    isValidated bool
}

func NewValidatedRecord(record Record) (*ValidatedRecord, error) {
    if err := record.validate(); err != nil {
        return nil, err
    }

    return &ValidatedRecord{
        Record: record,
        isValidated: true,
    }, nil
}

