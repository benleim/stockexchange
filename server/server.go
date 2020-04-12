package main

import (
	"context"
    "fmt"
    "log"
	"net/http"

	"github.com/labstack/echo"

	// "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// ----- Structs -----
type Order struct {
	Type   	string `json:"type"`
	Shares 	string `json:"int"`
}

func main() {
	e := echo.New()

	// Mongo setup
	setupMongo()

	// List endpoints
	e.GET("/stock/:ticker", getPrice)
	e.POST("/order", order)

	// Start http server
	e.Logger.Fatal(e.Start(":1323"))
}

func setupMongo() {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
	
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
	fmt.Println("Connected to MongoDB!")
}


// ----- Endpoints -----

// GET - Get Stock Price
func getPrice(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("ticker")
  	return c.String(http.StatusOK, id)
}

// POST - Order
func order(c echo.Context) (err error) {
	u := new(Order)
	if err = c.Bind(u); err != nil {
	  return
	}
	return c.JSON(http.StatusOK, u)
}