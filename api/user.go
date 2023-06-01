package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func User(c *fiber.Ctx) error {

	var tokenstring string
	const SecretKey = "secret"

	c.Request().Header.VisitAllCookie(func(key, value []byte) {
		tokenstring = string(value)
		fmt.Println("res cookieKey", string(key), "value", tokenstring)
	})

	// Use the tokenstring variable here
	fmt.Println("Tokenstring:", tokenstring)

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		// Token is invalid or expired
		// Handle the authentication failure
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var email string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email = claims["iss"].(string)
	}
	return c.SendString("Logged in User as " + email)

}
