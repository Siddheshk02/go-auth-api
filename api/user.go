package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Info struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func User(c *fiber.Ctx) error {

	var tokenstring string
	const SecretKey = "secret"

	c.Request().Header.VisitAllCookie(func(key, value []byte) {
		tokenstring = string(value)
	})

	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		// Token is invalid or expired
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

	var mongoUri = "mongodb+srv://Siddhesh:Sid42002@cluster0.puhgp.mongodb.net/?retryWrites=true&w=majority"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoUri,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been added to the access list. Error: ")
		panic(err)
	}

	userinfo := new(Info)

	if err := c.BodyParser(userinfo); err != nil {
		return err
	}

	userinfo.Email = email

	collection := client.Database("users-info").Collection("users-data")

	info := []interface{}{userinfo}

	insertManyResult, err := collection.InsertMany(context.TODO(), info)
	_ = insertManyResult
	if err != nil {
		fmt.Println("Something went wrong trying to insert the new document:")
		panic(err)
	}

	return c.SendString("Logged in as " + email + ".\nData added Successfully.")

}
