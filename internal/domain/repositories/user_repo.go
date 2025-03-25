package repositories


type UserRepository interface {
    Create(userID int64) error
}
