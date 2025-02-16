package payment

import (
	"fmt"
	"log/slog"

	"github.com/Razzle131/merchStore/internal/consts"
	"github.com/Razzle131/merchStore/internal/model"
	"github.com/Razzle131/merchStore/internal/repository/merchRepo"
	"github.com/Razzle131/merchStore/internal/repository/transactionRepo"
	"github.com/Razzle131/merchStore/internal/repository/userRepo"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

type PaymentService struct {
	userRepo  userRepo.UserRepoInterface
	merchRepo merchRepo.MerchRepoInterface
	transRepo transactionRepo.TransactionRepoInterface
}

func New(
	ur userRepo.UserRepoInterface,
	mr merchRepo.MerchRepoInterface,
	tr transactionRepo.TransactionRepoInterface,
) *PaymentService {
	return &PaymentService{
		userRepo:  ur,
		merchRepo: mr,
		transRepo: tr,
	}
}

func (s *PaymentService) SendCoin(to, from string, amount int) error {
	userFrom, err := s.userRepo.GetUserByLogin(from)
	if err != nil && err == serverErrors.ErrUserNotFound {
		return serverErrors.ErrUserNotFound
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	userTo, err := s.userRepo.GetUserByLogin(to)
	if err != nil && err == serverErrors.ErrUserNotFound {
		return serverErrors.ErrUserNotFound
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	err = userFrom.Wallet.AddOutTransaction(to, from, amount)
	if err != nil && err == serverErrors.ErrNotEnoughtMoney {
		return serverErrors.ErrNotEnoughtMoney
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	err = s.transRepo.AddTransaction(from, to, amount)
	if err != nil {
		return serverErrors.ErrInternal
	}

	err = s.userRepo.UpdateUser(userFrom)
	if err != nil && err == serverErrors.ErrNotAllowed {
		return serverErrors.ErrNotAllowed
	} else if err != nil {
		return serverErrors.ErrInternal
	}
	err = s.userRepo.UpdateUser(userTo)
	if err != nil && err == serverErrors.ErrNotAllowed {
		return serverErrors.ErrNotAllowed
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	slog.Debug(fmt.Sprint(s.userRepo))

	return nil
}

func (s *PaymentService) GetWalletInfo(userId int) (model.WalletInfo, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil && err == serverErrors.ErrUserNotFound {
		return model.WalletInfo{}, serverErrors.ErrUserNotFound
	} else if err != nil {
		return model.WalletInfo{}, serverErrors.ErrInternal
	}

	return user.Wallet, nil
}

func (s *PaymentService) BuyMerch(userId int, item string) error {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil && err == serverErrors.ErrUserNotFound {
		return serverErrors.ErrUserNotFound
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	price, err := s.merchRepo.GetMerchPrice(item)
	if err != nil && err == serverErrors.ErrItemNotFound {
		return serverErrors.ErrItemNotFound
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	slog.Debug("BuyMerch user " + fmt.Sprint(user))
	slog.Debug("BuyMerch price " + fmt.Sprint(price))

	err = user.Wallet.AddOutTransaction(consts.AvitoShop, user.Login, price)
	if err != nil && err == serverErrors.ErrNotEnoughtMoney {
		return serverErrors.ErrNotEnoughtMoney
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	err = user.Items.AddItem(item)
	if err != nil {
		return serverErrors.ErrInternal
	}

	err = s.userRepo.UpdateUser(user)
	if err != nil && err == serverErrors.ErrNotAllowed {
		return serverErrors.ErrNotAllowed
	} else if err != nil {
		return serverErrors.ErrInternal
	}

	return nil
}
