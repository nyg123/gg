//go:build wireinject

package dao

import (
	"github.com/google/wire"
)

func MakeUser() *User {
	wire.Build(
		NewDb,
		NewUser,
	)
	return &User{}
}
