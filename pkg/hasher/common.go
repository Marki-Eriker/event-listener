package hasher

import c "github.com/kaatinga/const-errs"

const ErrInvalidLength c.Error = "string is too long"

type Hasher interface {
	Hash(password string) (string, error)
	CheckHash(password string, hash string) bool
}
