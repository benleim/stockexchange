package bot

type Order struct {
	Side   string
	Qty    int
	Price  float64
	Open   bool
	Ticker string
}

func main() {
	tradingLoop()
}

func tradingLoop() {
	for {

	}
}
