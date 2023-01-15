package request

import (
	"context"
	c "github.com/kaatinga/const-errs"
	"github.com/marki-eriker/event-listener/entity/user"
	uuid "github.com/satori/go.uuid"
)

type ctxKey uint8

const (
	TraceID ctxKey = iota
	UserRole
)

const ErrRoleNotFoundInContext c.Error = "role not found in context"

// GetUserRole - Вернет роль пользователя из контекста
func GetUserRole(ctx context.Context) (user.Role, error) {
	v := ctx.Value(UserRole)

	if v == nil {
		return 0, ErrRoleNotFoundInContext
	}

	role, ok := v.(user.Role)
	if !ok {
		return 0, ErrRoleNotFoundInContext
	}

	return role, nil
}

// GetTraceID - вернет id из контекста или сгенерирует новый в случае отсутствия
func GetTraceID(ctx context.Context) string {
	v := ctx.Value(TraceID)
	if v == nil {
		return uuid.NewV4().String()
	}

	id, ok := v.(string)
	if !ok {
		return uuid.NewV4().String()
	}

	return id
}
