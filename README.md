# StockExchange
A mock stock exchange build in **Go** and (soon to be) deployed with **Kubernetes**.

Runs trading bot(s) that randomly place BUY or SELL orders with the exchange server. These orders are filled using different order matching algorithms (FIFO, Pro-Rata).

## Getting started
To run the stock exchange system simply run
```bash
docker-compose up
```
This will start the  `exchange`, `mongodb`, and `trading_bot` containers. Note that the trading bot places either a BUY or SELL order with a 50 / 50 chance. Naturally, the price of the stock will not fluctuate much.

Simply changing the supply / demand by changing the proportion of BUY or SELL trades is a great lesson of economics. More BUY orders will naturally increase the price. More SELL orders will decrease the price. I'm currently  working on randomly distributing this proportion.

## Server Request Schema
POST Order Body
```json
{
	"side": string,
	"qty": int,
	"price": float,
	"ticker": string
}
```
GET Stock Price
```bash 
GET http://exchange:1323/stock/:ticker
```

## Database Schema
Order
```json
{
	"_id": DocumentID,
	"side": bool,
	"time": ISO Timestamp,
	"qty": int,
	"price": float,
	"open": bool,
	"ticker": string,
	"remaining": int
}
```
Transaction
```json
{
	"to": string (bid _id),
	"from": string (ask _id),
	"Qty": int,
	"price": float,
	"ticker": string,
	"time": ISO Timestamp
}
```


## Completed Features
* Server
    * FIFO matching algo

* Basic Trading Bot

## Features in Development
* Server (*The Stock Exchange*)
    * Pro-Rata Algos
    * Short selling
    * Options
    * HFT Simulation (Front Running)
    * Multiple Exchanges

* Trading Bots (*The Traders*)
    * Risk adversion levels
    * Classes of Traders
    * Capital allocation

* Xchg Cli
