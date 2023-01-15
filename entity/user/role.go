package user

type Role uint16

const (
	Admin   Role = 100
	Analyst Role = 200
)

// CheckRole - проверит, существует ли такая роль.
// Может вернуть ErrInvalidRole
func checkRole(num uint16) (Role, error) {
	if Role(num) != Admin && Role(num) != Analyst {
		return 0, ErrInvalidRole
	}

	return Role(num), nil
}
