package processor

import c "github.com/kaatinga/const-errs"

const (
	ErrUnexpectedDatabaseBehavior c.Error = "unexpected database behavior"
	ErrUserNotFound               c.Error = "user not found"
	ErrUserAlreadyExists          c.Error = "user already exists"
	ErrInvalidCredentials         c.Error = "invalid credentials"
	ErrTokenGeneration            c.Error = "token generation error"
	ErrNotVerified                c.Error = "user not verified"
)
