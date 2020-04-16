db.auth('admin-user', 'admin-password')

// Listed Stocks (Collection)
db.createCollection("stocks", { capped : true, size : 5242880, max : 5000 } )
db.stocks.insert( { ticker: "AAPL", name: "Apple", IPO: 15.50 } )
db.stocks.insert( { ticker: "MSFT", name: "Microsoft", IPO: 20.25 } )

// Orders (Collection)
db.createCollection("orders", { capped : true, sizee : 5252880, max : 5000 } )
db.orders.insert( { id: 0, side: "SELL", time: "2020-04-15T16:30:17.134221-04:00", qty: 100, price: 15.40, open: true, filledWith: [] } )
db.orders.insert( { id: 0, side: "SELL", time: "2020-04-15T16:30:17.134221-04:01", qty: 80, price: 15.39, open: true, filledWith: [] } )
db.orders.insert( { id: 0, side: "SELL", time: "2020-04-15T16:30:17.134221-04:02", qty: 70, price: 15.38, open: true, filledWith: [] } )