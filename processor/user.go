package processor

import (
	"context"
	"github.com/marki-eriker/event-listener/db"
	"github.com/marki-eriker/event-listener/entity/user"
	"github.com/marki-eriker/event-listener/pkg/recovery"
	"github.com/marki-eriker/event-listener/pkg/request"
)

// RegisterUser - зарегистрирует нового пользователя.
// Может вернуть user.ErrLowLoginLength, user.ErrLowPasswordComplexity, user.ErrInvalidRole,
// ErrUserAlreadyExists, ErrUnexpectedDatabaseBehavior
func (p *Processor) RegisterUser(ctx context.Context, login, password string, role uint16) error {
	defer recovery.Recover("processor.RegisterUser")

	ll := p.Log(request.GetTraceID(ctx)).WithField("login", login).WithField("role", role)

	ll.Debug("RegisterUser begin")
	defer ll.Debug("RegisterUser end")

	u, err := user.New(login, password, role, p.hasher)
	if err != nil {
		ll.Errorf("unable to create valid user struct: %s", err)
		return err
	}

	err = p.store.User.Insert(ctx, u)
	if err != nil && err != db.ErrDuplicateKey {
		ll.Errorf("unable to save user: %s", err)
		return ErrUnexpectedDatabaseBehavior
	}

	if err == db.ErrDuplicateKey {
		ll.Debug("user already exists")
		return ErrUserAlreadyExists
	}

	return nil
}

// VerifyUser - подтвердит регистрацию пользователя.
// Может вернуть ErrUnexpectedDatabaseBehavior, ErrUserNotFound, user.ErrAlreadyVerified
func (p *Processor) VerifyUser(ctx context.Context, id uint) error {
	defer recovery.Recover("processor.VerifyUser")

	ll := p.Log(request.GetTraceID(ctx)).WithField("id", id)

	ll.Debug("VerifyUser begin")
	defer ll.Debug("VerifyUser end")

	u, err := p.store.User.GetByID(ctx, id)
	if err != nil && err != db.ErrRecordNotFound {
		ll.Errorf("unable to find user: %s", err)
		return ErrUnexpectedDatabaseBehavior
	}

	if err == db.ErrRecordNotFound {
		ll.Debug("user not found")
		return ErrUserNotFound
	}

	err = u.Verify()
	if err != nil {
		ll.Debug("user already verified")
		return err
	}

	err = p.store.User.UpdateVerify(ctx, u.ID, true)
	if err != nil {
		ll.Errorf("unable to update user: %s", err)
		return ErrUnexpectedDatabaseBehavior
	}

	return nil
}

// LoginUser - вернет авторизационный токен подтвержденному пользователю.
// Может вернуть ErrUnexpectedDatabaseBehavior, ErrInvalidCredentials, ErrTokenGeneration
func (p *Processor) LoginUser(ctx context.Context, login, password string) (string, error) {
	defer recovery.Recover("processor.LoginUser")

	ll := p.Log(request.GetTraceID(ctx))

	ll.Debug("LoginUser begin")
	defer ll.Debug("LoginUser end")

	u, err := p.store.User.GetByLogin(ctx, login)
	if err != nil && err != db.ErrRecordNotFound {
		ll.Errorf("unable to find user: %s", err)
		return "", ErrUnexpectedDatabaseBehavior
	}

	if err == db.ErrRecordNotFound {
		ll.Debug("user not found")
		return "", ErrInvalidCredentials
	}

	ok := p.hasher.CheckHash(password, u.Password)
	if !ok {
		ll.Debug("password mismatch")
		return "", ErrInvalidCredentials
	}

	if !u.Verified {
		ll.Debug("user not verified")
		return "", ErrNotVerified
	}

	token, err := p.token.Generate(u.Role)
	if err != nil {
		ll.Errorf("unable to generate token: %s", err)
		return "", ErrTokenGeneration
	}

	return token, nil
}
