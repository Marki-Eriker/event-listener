package hasher

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	cost int
}

func NewBcryptHasher(cost int) *Bcrypt {
	switch {
	case cost == 0:
		cost = bcrypt.DefaultCost
	case cost < bcrypt.MinCost:
		cost = bcrypt.MinCost
	case cost > bcrypt.MaxCost:
		cost = bcrypt.MaxCost
	}

	return &Bcrypt{cost: cost}
}

func (bc *Bcrypt) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bc.cost)
	return string(b), err
}

func (bc *Bcrypt) CheckHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
