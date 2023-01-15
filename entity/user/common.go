package user

import c "github.com/kaatinga/const-errs"

const (
	ErrAlreadyVerified       c.Error = "user already verified"
	ErrLowPasswordComplexity c.Error = "low password complexity"
	ErrLowLoginLength        c.Error = "low login length"
	ErrInvalidRole           c.Error = "invalid role value"
)
