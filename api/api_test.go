package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuth(t *testing.T) {
	app := fiber.New()

	app.Post("/auth", func(c *fiber.Ctx) error {
		userdata := new(Data)

		userdata.Email = "Sid@abc.com"
		userdata.Password = "123456789"

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

		collection := client.Database("users-info").Collection("users")

		ch, ct := Check(userdata, client, collection)

		if ct == 0 {

			insertOneResult, err := collection.InsertOne(context.TODO(), userdata)
			_ = insertOneResult
			if err != nil {
				fmt.Println("Something went wrong trying to insert the new document:")
				panic(err)
			}
		} else if ct == 2 {
			return c.SendString(ch)
		}

		const SecretKey = "secret"

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    userdata.Email,
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		response := map[string]string{"token": token}
		return c.JSON(response)
	})
}

func TestUser(t *testing.T) {
	app := fiber.New()

	app.Post("/user", func(c *fiber.Ctx) error {
		userinfo := new(Info)

		userinfo.Name = "Siddhesh"
		userinfo.Email = "Sid@abc.com"
		userinfo.Phone = "123456789"

		var tokenstring string
		const SecretKey = "secret"

		c.Request().Header.VisitAllCookie(func(key, value []byte) {
			tokenstring = string(value)
		})

		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {

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

		userinfo.Email = email

		collection := client.Database("users-info").Collection("users-data")

		info := []interface{}{userinfo}

		insertManyResult, err := collection.InsertMany(context.TODO(), info)
		_ = insertManyResult
		if err != nil {
			fmt.Println("Something went wrong trying to insert the new document:")
			panic(err)
		}

		// return c.SendString("Logged in as " + email + ".\nData added Successfully.")

		response := map[string]string{"message": "User information stored successfully"}
		return c.JSON(response)
	})

}

func TestMain(m *testing.M) {
	m.Run()
}
