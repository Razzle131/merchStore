package authorization

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Razzle131/merchStore/internal/model"
	"github.com/Razzle131/merchStore/internal/repository/userRepo"
	"github.com/Razzle131/merchStore/internal/serverErrors"
	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

const (
	jwtExpire = 10 * time.Minute
)

type jwtClaims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type AuthorizationService struct {
	userRepo userRepo.UserRepoInterface
}

func New(ur userRepo.UserRepoInterface) *AuthorizationService {
	return &AuthorizationService{
		userRepo: ur,
	}
}

func (s *AuthorizationService) AuthenticateUser(login, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		login,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	user, err := s.userRepo.GetUserByLogin(login)
	if err != nil && err == serverErrors.ErrUserNotFound {
		err = s.userRepo.AddUser(login, password)
		if err != nil {
			return "", serverErrors.ErrInternal
		}

		return token.SignedString(jwtSecret)
	} else if err != nil {
		return "", serverErrors.ErrInternal
	}

	if user.Password != password {
		return "", serverErrors.ErrBadCreditonals
	}

	return token.SignedString(jwtSecret)
}

func (s *AuthorizationService) AuthorizeUser(token string) (model.User, error) {
	splittedToken := strings.Split(token, " ")
	if len(splittedToken) < 2 {
		slog.Debug("bad token")
		return model.User{}, serverErrors.ErrBadToken
	}

	tmp := jwtClaims{}
	parsedToken, _ := jwt.ParseWithClaims(splittedToken[1], &tmp, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	slog.Debug("AuthorizeUser: " + fmt.Sprint(tmp.Login))

	if parsedToken.Valid {
		slog.Debug("token valid")
		return s.userRepo.GetUserByLogin(tmp.Login)
	}

	return model.User{}, serverErrors.ErrBadToken
}
