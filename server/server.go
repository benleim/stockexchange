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

type Transaction struct {
	To     string
	From   string
	Qty    int
	Price  float64
	Ticker string
	Time   time.Time
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
	findOptions.SetSort(bson.D{
		{"price", priceSort},
		{"time", -1},
	})
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

// Returns documentID string
func writeOrder(cli *mongo.Client, u *Order) string {
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

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex()
	} else {
		return "error"
	}
}

func writeTransaction(cli *mongo.Client, trans *Transaction, otherID string) {
	col := cli.Database("exchange").Collection("transactions")
	// Add missing fields to transaction
	trans.Time = time.Now()
	if trans.From == "" {
		trans.From = otherID
	} else {
		trans.To = otherID
	}

	// Change to bsonM
	transDoc := bson.M{
		"to":     trans.To,
		"from":   trans.From,
		"qty":    trans.Qty,
		"price":  trans.Price,
		"ticker": trans.Ticker,
		"time":   trans.Time,
	}

	result, insertErr := col.InsertOne(context.TODO(), transDoc)
	if insertErr != nil {
		println("Transaction Insert ERROR:", insertErr)
	} else {
		println("Transaction Insert Result:", result)
	}
}

func orderToTrans(order *Order, qty int) Transaction {
	trans := Transaction{
		Qty:    qty,
		Price:  order.Price,
		Ticker: order.Ticker,
	}
	if order.Side == "BUY" {
		trans.To = order.Id.Hex()
	} else {
		trans.From = order.Id.Hex()
	}

	return trans
}

// ----- Endpoints -----

// GET - Get Stock IPO Price
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

	sum := 0
	queue := make([]Transaction, 0)
	cur := getOrders(client, u.Ticker, u.Price, side)
	for cur.Next(nil) {
		// Decode order
		var result Order
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		println(strconv.FormatFloat(result.Price, 'f', 2, 64), result.Qty, primitive.ObjectID.Hex(result.Id))
		sum += result.Remaining
		var amount int
		if sum <= u.Qty {
			amount = result.Remaining
		} else {
			amount = result.Remaining - (sum - u.Qty)
			println("remaining", amount)
		}
		updateOrder(client, result, amount)
		trans := orderToTrans(&result, amount)
		queue = append(queue, trans)

		// Stop when order is filled
		if sum >= u.Qty {
			break
		}
	}

	// Determine amt remaining
	if sum < u.Qty {
		u.Remaining = (u.Qty - sum)
	} else {
		u.Remaining = 0
	}

	// Write new order to db
	newDocID := writeOrder(client, u)
	println("New DocumentID", newDocID)

	// Log transactions
	for len(queue) > 0 {
		writeTransaction(client, &queue[0], newDocID)
		queue = queue[1:]
	}

	return c.JSON(http.StatusOK, u)
}
