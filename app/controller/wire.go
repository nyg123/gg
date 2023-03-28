//go:build wireinject

package controller

import (
	"gg/app/dao"
	"github.com/google/wire"
)

func MakeUser() *User {
	wire.Build(
		dao.MakeUser,
		NewUser,
	)

	return &User{}
}
