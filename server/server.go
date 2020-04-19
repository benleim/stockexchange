package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ----- Structs -----

type Order struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Side      string             `json:"side"`
	Time      time.Time          `json:"time"`
	Qty       int                `json:"qty"`
	Price     float64            `json:"price"`
	Open      bool               `json:"open"`
	Ticker    string             `json:"ticker"`
	Remaining int                `json:"remaining"`
}

type OrderW struct {
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

// ----- Utils -----
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

func getOrders(cli *mongo.Client, ticker string, maxPrice float64, side string) *mongo.Cursor {
	// Side vars for ordering
	priceSort := 1
	maxQuery := "$lte"
	if side == "BUY" {
		priceSort = -1
		maxQuery = "$gte"
	}

	collection := cli.Database("exchange").Collection("orders")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"price", priceSort}})
	query := bson.D{
		{"side", side},
		{"open", true},
		{"ticker", ticker},
		{"price", bson.D{
			{maxQuery, maxPrice},
		}},
	}
	cur, err := collection.Find(context.TODO(), query, findOptions)
	println("Fetched order data from mongo")
	if err != nil {
		log.Fatal(err)
	}

	return cur
}

func updateOrder(cli *mongo.Client, o Order, amount int) {
	col := cli.Database("exchange").Collection("orders")
	filter := bson.M{"_id": o.Id}
	remaining := o.Remaining - amount
	open := remaining > 0
	update := bson.M{"$set": bson.M{
		"open":      open,
		"remaining": remaining,
	}}
	result, err := col.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		log.Fatal(err)
	}
	println("Update Result: ", result)
}

func writeOrder(cli *mongo.Client, u *Order) {
	orderDoc := bson.M{
		"side":      u.Side,
		"time":      time.Now(),
		"qty":       u.Qty,
		"price":     u.Price,
		"open":      u.Remaining > 0,
		"ticker":    u.Ticker,
		"remaining": u.Remaining,
	}
	col := cli.Database("exchange").Collection("orders")
	result, insertErr := col.InsertOne(context.TODO(), orderDoc)
	if insertErr != nil {
		println("Order Insert ERROR:", insertErr)
	} else {
		println("Order Insert Result:", result)
	}
}

func writeTransaction(cli *mongo.Client) {

}

// ----- Endpoints -----

// GET - Get Stock Price
func getPrice(c echo.Context) error {
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

	// Get mongo client
	client := getMongoClient()
	// Get opposite orders
	side := "SELL"
	if u.Side == "SELL" {
		side = "BUY"
	}
	cur := getOrders(client, u.Ticker, u.Price, side)
	sum := 0
	for cur.Next(nil) {
		var result Order
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		println(strconv.FormatFloat(result.Price, 'f', 2, 64), result.Qty, primitive.ObjectID.Hex(result.Id))
		sum += result.Remaining
		if sum <= u.Qty {
			updateOrder(client, result, result.Remaining)
		} else {
			remaining := result.Remaining - (sum - u.Qty)
			println("remaining", remaining)
			updateOrder(client, result, remaining)
		}

		if sum >= u.Qty {
			break
		}
	}

	// Write order to DB
	if sum < u.Qty {
		u.Remaining = (u.Qty - sum)
	} else {
		u.Remaining = 0
	}
	writeOrder(client, u)

	return c.JSON(http.StatusOK, u)
}
