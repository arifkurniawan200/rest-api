package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
	"template/internal/model"
	"template/internal/repository"
	"template/internal/utils"
)

type UserHandler struct {
	u repository.UserRepository
	t repository.TransactionRepository
}

func NewUserUsecase(u repository.UserRepository, t repository.TransactionRepository) UserUcase {
	return &UserHandler{u, t}
}

func (u UserHandler) GetUserInfoByEmail(ctx echo.Context, email string) (model.User, error) {
	var (
		err error
	)
	userInfo, err := u.u.GetUserByEmail(email)
	return userInfo, err
}

func (u UserHandler) RegisterCustomer(ctx echo.Context, c model.RequestRegisterUser) error {
	var (
		err error
		g   errgroup.Group
	)

	g.Go(func() error {
		// hash password
		c.Password, err = utils.HashPassword(c.Password)
		if err != nil {
			log.Errorf("error when hash password ")
			return err
		}
		return err
	})

	if err = g.Wait(); err != nil {
		return err
	}

	err = u.u.RegisterUser(c)
	if err != nil {
		log.Errorf("[usecase][RegisterCustomer] error when RegisterUser: %s", err.Error())
		return err
	}

	return nil
}
