package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rehandwi03/test-case-backend-majoo/internal/helper"
	"log"
	"os"
	"strings"
)

type TokenDetail struct {
	UserID uuid.UUID `json:"user_id"`
}

func exractToken(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("token not found")
	}
	token := strings.Split(bearerToken, " ")
	if len(token) == 2 {
		return token[1], nil
	}
	return "", errors.New("token not found")
}

// verfiy token format
func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString, err := exractToken(c)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("APP_SECRET")), nil
		},
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(c *fiber.Ctx) (*TokenDetail, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}

		userIDUUID, err := uuid.Parse(userId)
		if err != nil {
			return nil, err
		}

		return &TokenDetail{
			UserID: userIDUUID,
		}, nil
	}
	return nil, err
}

func checkAuthToken(c *fiber.Ctx) (*TokenDetail, error) {
	extracTokenDetails, err := ExtractTokenMetadata(c)
	if err != nil {
		return nil, err
	}

	return extracTokenDetails, nil
}

func JwtProtected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenDetails, err := checkAuthToken(ctx)
		if err != nil {
			log.Printf("checkAuthToken mdw :" + err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				helper.ErrorResponse{
					Status:  "failed",
					Message: "StatusUnauthorized",
					Errors:  err.Error(),
				},
			)
		}

		ctx.Locals("user_id", tokenDetails.UserID)
		return ctx.Next()
	}
}
