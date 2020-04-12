package main

import (
	"net/http"
	"github.com/labstack/echo"
)

// ----- Structs -----
type Order struct {
	Type   	string `json:"type"`
	Shares 	string `json:"int"`
}

func main() {
	e := echo.New()

	e.GET("/stock/:ticker", getPrice)
	e.POST("/order", order)

	e.Logger.Fatal(e.Start(":1323"))
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