package request

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
)

// Value - тип строки, с возможностью его преобразования в иные значения
type Value string

// MustUInt - преобразовывает значение в целое положительное число
func (rv Value) MustUInt() uint {
	res, err := strconv.Atoi(string(rv))
	if err != nil {
		return 0
	}

	if res < 1 {
		return 0
	}

	return uint(res)
}

// String - возвращает строковое значение
func (rv Value) String() string {
	return string(rv)
}

func (rv Value) MustUUID() *uuid.UUID {
	res, err := uuid.FromString(string(rv))
	if err != nil {
		return nil
	}

	return &res
}
