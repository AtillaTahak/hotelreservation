package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWTAuthentication")
	token, ok := c.GetReqHeaders()["X-Access-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if err := parseToken(token); err != nil {
		return err
	}
	fmt.Println(token)
	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println(secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parser JWT toke", err)
		return fmt.Errorf("unauthorized")
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unauthorized")
	
}
