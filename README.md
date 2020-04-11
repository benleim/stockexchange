# StockExchange
A mock stock exchange build in **Go** and deployed in **Kubernetes**.

Runs trading bots that randomly place BUY or SELL orders with the exchange Django API. These orders are filled using different matching algorithms (FIFO, Pro-Rata). The order matching algorithm can be specified by the `xchg` cli.

The trading bots will randomly place BUY or SELL orders with a normal distribution of prices centered at the current lowest bid price and highest ask price at the time.

## Completed Features

## Features in Development
* API (*The Stock Exchange*)
    * FIFO, Pro-Rata Algos
    * Short selling
    * Options
    * HFT Simulation (Front Running)
    * Multiple Exchanges

* Trading Bots (*The Traders*)
    * Risk adversion levels
    * Classes of Traders
    * Capital allocation

* Xchg Cli
