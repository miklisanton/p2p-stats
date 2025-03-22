package entities

type ValidatedUser struct {
    User
    isValidated bool
}

func NewValidatedUser(user User) (*ValidatedUser, error) {
    if err := user.validate(); err != nil {
        return nil, err
    }

    return &ValidatedUser{
        User: user,
        isValidated: true,
    }, nil
}

