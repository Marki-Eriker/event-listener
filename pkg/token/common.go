package token

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/marki-eriker/event-listener/entity/user"
	"time"
)

type Generator interface {
	Generate(role user.Role) (string, error)
	Parse(string) (*jwt.Token, error)
}

type Options struct {
	Secret string
	Ttl    time.Duration
}
