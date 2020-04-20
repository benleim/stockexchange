db.auth('admin-user', 'admin-password')

// Listed Stocks (Collection)
db.createCollection("stocks", { capped : true, size : 5242880, max : 5000 } )
db.stocks.insert( { ticker: "AAPL", name: "Apple", IPO: 15.50 } )
db.stocks.insert( { ticker: "MSFT", name: "Microsoft", IPO: 20.25 } )

// Orders (Collection)
db.createCollection("orders", { capped : false } )

// Transactions (Collection)
db.createCollection("transactions", { capped : false })
