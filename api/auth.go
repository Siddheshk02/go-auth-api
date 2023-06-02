package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Data struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Auth(c *fiber.Ctx) error {
	userdata := new(Data)

	if err := c.BodyParser(userdata); err != nil {
		return err
	}

	var mongoUri = "mongodb+srv://@cluster0.puhgp.mongodb.net/?retryWrites=true&w=majority"

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
		info := []interface{}{userdata}

		insertManyResult, err := collection.InsertMany(context.TODO(), info)
		_ = insertManyResult
		if err != nil {
			fmt.Println("Something went wrong trying to insert the new documents:")
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

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Minute * 30),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.SendString(ch + ". \nToken Generated Successfully.")

}

func Check(userdata *Data, client *mongo.Client, collection *mongo.Collection) (string, int) {

	var result Data
	var myFilter = bson.D{{Key: "email", Value: userdata.Email}}
	e := collection.FindOne(context.TODO(), myFilter).Decode(&result)
	if e != nil {
		return "New User Added", 0
	} else {
		if result.Password == userdata.Password {
			return "User Already Exists", 1
		} else {
			return "User Already Exists. [Incorrect Password]", 2
		}
	}

}
