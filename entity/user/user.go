package user

import (
	"github.com/marki-eriker/event-listener/pkg/hasher"
)

// User - Основная структура пользователя
type User struct {
	tableName struct{} `pg:"users"`

	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
	Verified bool   `json:"verified"`
}

// New - создаст нового пользователя готового к записи в БД.
// Проверит логин, пароль и роль на соответствие правилам.
// Захеширует пароль.
// Может вернуть ErrLowLoginLength, ErrLowPasswordComplexity, ErrInvalidRole
func New(login, password string, roleNum uint16, hasher hasher.Hasher) (*User, error) {

	err := checkCredentials(login, password)
	if err != nil {
		return nil, err
	}

	role, err := checkRole(roleNum)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Login:    login,
		Password: hashedPassword,
		Role:     role,
	}, nil
}

// Verify - отметит подтверждение регистрации пользователя.
// Может вернуть ErrAlreadyVerified
func (u *User) Verify() error {
	if u.Verified {
		return ErrAlreadyVerified
	}

	u.Verified = true

	return nil
}

func checkCredentials(login, password string) error {
	if len(login) < MinimumLoginLength {
		return ErrLowLoginLength
	}

	if len(password) < MinimumPasswordLength {
		return ErrLowPasswordComplexity
	}

	return nil
}
