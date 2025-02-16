package userRepo

import (
	"github.com/Razzle131/merchStore/internal/model"
)

type UserRepoInterface interface {
	AddUser(login, password string) error
	GetUserById(id int) (model.User, error)
	GetUserByLogin(login string) (model.User, error)
	UpdateUser(newUser model.User) error
}
