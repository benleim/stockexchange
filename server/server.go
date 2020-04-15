package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ----- Structs -----
type Order struct {
	Type   string `json:"type"`
	Shares string `json:"int"`
}

type Stock struct {
	Name   string  `json:"name"`
	Ticker string  `json:"ticker"`
	IPO    float64 `json:"ipo"`
}

func main() {
	e := echo.New()

	// List endpoints
	e.GET("/stock/:ticker", getPrice)
	e.POST("/order", order)

	// Start http server
	e.Logger.Fatal(e.Start(":1323"))
}

func getMongoClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017").
		SetAuth(options.Credential{
			AuthSource: "admin",
			Username:   "admin-user",
			Password:   "admin-password",
		})

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// ----- Endpoints -----

// GET - Get Stock Price
func getPrice(c echo.Context) error {
	// User ID from path `users/:id`
	ticker := c.Param("ticker")

	client := getMongoClient()

	// Reference "stocks" collection
	collection := client.Database("exchange").Collection("stocks")
	var result Stock
	err := collection.FindOne(context.TODO(), bson.D{{"ticker", ticker}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	ipo := strconv.FormatFloat(result.IPO, 'f', 2, 64)
	return c.String(http.StatusOK, ipo)
}

// POST - Order
func order(c echo.Context) (err error) {
	u := new(Order)
	if err = c.Bind(u); err != nil {
		return
	}
	return c.JSON(http.StatusOK, u)
}
