package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthorizeUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User Unauthorized: missing Authorization Header"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := VerifyToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User Unauthorized: " + err.Error(),
		})
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token: missing user ID"})
	}
	c.Locals("userId", userId)

	return c.Next()
}

func CreateToken(UserId string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	byteKey := []byte(secretKey)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    UserId,
		"issuer": "blogpost",
		"Role":   "User",
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	})

	Token, err := claims.SignedString(byteKey)
	if err != nil {
		return "", err
	}
	return Token, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	byteKey := []byte(secretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		fmt.Println(byteKey, token)
		return byteKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token is expired")
		}
	}

	return claims, nil
}
