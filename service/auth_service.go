package service

import (
	"context"
	"github.com/golang-jwt/jwt"
	custom_error "github.com/rehandwi03/test-case-backend-majoo/internal/error"
	"github.com/rehandwi03/test-case-backend-majoo/repository"
	"github.com/rehandwi03/test-case-backend-majoo/request"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, request *request.LoginRequest) (map[string]interface{}, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepo: userRepository}
}

func (a *authService) Login(ctx context.Context, request *request.LoginRequest) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"where": map[string]interface{}{
			"default": map[string]interface{}{
				"email = ?": request.Email,
			},
		},
	}

	checkUser, err := a.userRepo.GetByParam(ctx, params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &custom_error.BadRequest{"email or password is incorrect"}
		}

		return nil, err
	}

	ok, err := checkUser.ComparePassword(request.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !ok {
		return nil, &custom_error.BadRequest{Message: "email or password is incorrect"}
	}

	token, err := a.GenerateToken(ctx, checkUser.ID.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resToken := map[string]interface{}{
		"access_token": token,
	}

	return resToken, nil
}

func (a *authService) GenerateToken(
	ctx context.Context, userId string,
) (token string, err error) {
	exp := time.Now().Add(time.Hour * 15).Unix()

	// create access token
	acessTokenClaims := jwt.MapClaims{}
	acessTokenClaims["authorized"] = true
	acessTokenClaims["user_id"] = userId
	acessTokenClaims["exp"] = exp

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, acessTokenClaims)
	// assign access token to tokendetails struct
	token, err = accessToken.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return token, err
	}

	return token, err
}
