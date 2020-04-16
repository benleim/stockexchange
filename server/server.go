package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ----- Structs -----

type Order struct {
	Id    int       `json:"id"`
	Side  string    `json:"side"`
	Time  time.Time `json:"time"`
	Qty   int       `json:"qty"`
	Price float64   `json:"price"`
	Open  bool      `json:"open"`
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

func getSellOrders(cli *mongo.Client, maxPrice float64) *mongo.Cursor {
	collection := cli.Database("exchange").Collection("orders")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"price", -1}})
	query := bson.D{
		{"price", bson.D{
			{"$lte", maxPrice},
		}},
	}
	cur, err := collection.Find(context.TODO(), query, findOptions)
	println("Fetched SELL data from mongo")
	if err != nil {
		log.Fatal(err)
	}

	// defer cur.Close(nil)
	return cur
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
	// Decode order req body
	u := new(Order)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Time = time.Now()

	client := getMongoClient()
	cur := getSellOrders(client, u.Price)
	for cur.Next(nil) {
		var result Order
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		println(strconv.FormatFloat(result.Price, 'f', 2, 64))
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, u)
}
