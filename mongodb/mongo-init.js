db.auth('admin-user', 'admin-password')

// Listed Stocks (Collection)
db.createCollection("stocks", { capped : true, size : 5242880, max : 5000 } )
db.stocks.insert( { ticker: "AAPL", name: "Apple", IPO: 15.50 } )
db.stocks.insert( { ticker: "MSFT", name: "Microsoft", IPO: 20.25 } )

// Orders (Collection)
db.createCollection("orders", { capped : false } )
db.orders.insert( { ticker: "AAPL", side: "SELL", time: ISODate("2020-04-20T01:51:37.829Z"), qty: 100, price: 15.40, open: true, remaining: 100 } )
db.orders.insert( { ticker: "AAPL", side: "SELL", time: ISODate("2020-04-21T01:51:37.829Z"), qty: 80, price: 15.39, open: true, remaining: 80 } )
db.orders.insert( { ticker: "AAPL", side: "SELL", time: ISODate("2020-04-19T01:51:37.829Z"), qty: 70, price: 15.38, open: true, remaining: 70 } )

// Transactions (Collection)
db.createCollection("transactions", { capped : false })
db.transactions.insert( { to: "", from: "", ticker: "AAPL", qty: 10, price: 12.50, time: ISODate("2020-04-15T01:51:37.829Z") } )