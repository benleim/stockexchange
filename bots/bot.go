package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type Order struct {
	Side   string
	Qty    int
	Price  float64
	Open   bool
	Ticker string
}

type Price struct {
	Price float64 `json:"price"`
}

func main() {
	tradingLoop()
}

func tradingLoop() {
	for {
		// Set random timeout
		timeout()

		// Get Price
		price := getPrice("AAPL")

		// Random number of shares (10 - 200)
		shares := rand.Intn(200-10) + 10
		println("Shares: ", shares)

		// Random price deviation from BUY (0.01 - 0.10 above or below)
		dev := -0.05 + rand.Float64()*(0.05-(-0.05))
		newPrice := math.Round((price+dev)*100) / 100
		fmt.Printf("New Price: %f\n\n", newPrice)

		// BUY or SELL ?
		buyorsell := rand.Float64()
		var side string
		if buyorsell < 0.5 {
			side = "BUY"
		} else {
			side = "SELL"
		}
		putOrder("AAPL", shares, newPrice, side)
	}
}

func timeout() {
	r := rand.Intn(10000)
	time.Sleep(time.Duration(r) * time.Microsecond)
}

func putOrder(ticker string, shares int, price float64, side string) {
	body := map[string]interface{}{
		"side":   side,
		"qty":    shares,
		"price":  price,
		"ticker": ticker,
	}

	bodyRep, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post("http://exchange:1323/order",
		"application/json", bytes.NewBuffer(bodyRep))
	if err != nil {
		print(err)
	}
	println(resp)
}

func getPrice(ticker string) float64 {
	var p Price
	response, err := http.Get("http://exchange:1323/stock/AAPL")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// read the payload, in this case, Jhon's info
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TRADING PRICE: %f\n", p.Price)

	// this is where the magic happens, I pass a pointer of type Person and Go'll do the rest
	err = json.Unmarshal(body, &p)
	return p.Price
}
