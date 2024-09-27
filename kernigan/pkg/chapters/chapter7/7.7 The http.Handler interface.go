package chapter7

import (
	"fmt"
	. "golang/pkg/chapters/chapter1/c_parse"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	fileName = "goods.html"
)

var (
	goodsTemplate = template.Must(template.New("goods").Parse(
		`
<!DOCTYPE>
<html>

<head>
	<title>My Store</title>
	<style>
		h1 {
			text-align: center;
		}
		table {
			border-collapse: collapse;
			width: 100%;
		}
		th, td {
			border: 1px solid black;
			padding: 8px;
		}
		td {
			text-align: left;
		}
		th {
			background-color: #dcfff2;
			text-align: center;
		}
	</style>
</head>

<body>
	<h1>Our Goods</h1>
	<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Price</th>
				</tr>
			</thead>
			
			<tbody>
				{{range $key, $value := .}}
					<tr>
						<td>{{$key}}</td>
						<td>{{$value}}</td>
					</tr>
				{{end}}
			</tbody>		
	</table>
</body>

</html>
`))
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.3g", d)
}

type database map[string]dollars

func startShopServer() {
	// Fill the database with some data
	db := database{"shoes": 100, "jeans": 30, "socks": 10}

	// Get a multiplexer to process several url's paths separately.
	multiplexer := http.NewServeMux()

	// Add the methods to multiplexer directly with .HandleFunc()
	multiplexer.HandleFunc("/list", db.list)
	multiplexer.HandleFunc("/price", db.price)
	multiplexer.HandleFunc("/create", db.create)
	multiplexer.HandleFunc("/delete", db.delete)
	multiplexer.HandleFunc("/update", db.update)

	// Pass the multiplexer instead of db
	log.Fatal(http.ListenAndServe("localhost:8080", multiplexer))
}

func UseShopServer() {
	time.AfterFunc(time.Second*1, func() {
		// Get all the items from the shop inventory
		Fetch("http://localhost:8080/list")
	})

	startShopServer()
}

/*
Provides the client with table of current goods in the store.
*/
func (db database) list(w http.ResponseWriter, req *http.Request) {
	// Create a file to store HTML table at runtime
	goodsFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("creating file; %s", err), http.StatusInternalServerError)
		return
	}

	// Try to parse data from the database to the HTML template. If everything is okay write filled data into the file
	err = goodsTemplate.Execute(goodsFile, db)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsing template; %s", err), http.StatusInternalServerError)
		return
	}
	goodsFile.Close()

	// Send the file to a client side
	http.ServeFile(w, req, fileName)
}

/*
Returns the price of an item if it's placed in inventory
*/
func (db database) price(w http.ResponseWriter, req *http.Request) {
	isPlaced := db.isGoodsPlaced(req.URL.Query().Get("item"))
	if !isPlaced {
		// Set the header of a response
		w.WriteHeader(http.StatusNotFound)
		// Write a message to a client side
		fmt.Fprintf(w, "getting the \"%s\" price from the database; no required item in the database\n", req.URL.Query().Get("item"))
		return
	}
	fmt.Fprintf(w, "%s\n", db[req.URL.Query().Get("item")])
}

/*
Puts an item into the shops's inventory
*/
func (db database) create(w http.ResponseWriter, req *http.Request) {
	if err := paramsValidation(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := db.isItemValid(req.URL.Query().Get("item"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if db.isGoodsPlaced(item) {
		http.Error(w, fmt.Sprintf("adding the item; item \"%s\" is already placed", req.URL.Query().Get("item")), http.StatusBadRequest)
		return
	}

	price, err := db.isPriceValid(req.URL.Query().Get("price"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db[req.URL.Query().Get("item")] = dollars(price)
}

/*
Removes an item from the database
*/
func (db database) delete(w http.ResponseWriter, req *http.Request) {
	isPlaced := db.isGoodsPlaced(req.URL.Query().Get("item"))
	if !isPlaced {
		// Set the header of a response
		w.WriteHeader(http.StatusNotFound)
		// Write a message to a client side
		fmt.Fprintf(w, "deleting the \"%s\" from the database; no required item in the database\n", req.URL.Query().Get("item"))
		return
	}
	delete(db, req.URL.Query().Get("item"))
}

/*
Changes the item's price
*/
func (db database) update(w http.ResponseWriter, req *http.Request) {
	if err := paramsValidation(req.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !db.isGoodsPlaced(req.URL.Query().Get("item")) {
		http.Error(w, fmt.Sprintf("finding the item; no item \"%s\" in the database ", req.URL.Query().Get("item")), http.StatusBadRequest)
		return
	}

	price, err := db.isPriceValid(req.URL.Query().Get("price"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db[req.URL.Query().Get("item")] = dollars(price)
}

/*
Checks the validity of the request's params.
*/
func paramsValidation(params url.Values) error {
	var (
		isItem, isPrice = false, false
	)
	if len(params) == 2 {
		for key := range params {
			if key == "item" {
				isItem = true
			}
			if key == "price" {
				isPrice = true
			}
		}

		if isItem && isPrice {
			return nil
		}

		return fmt.Errorf("while parsing %v; invalid params names", params)
	}

	return fmt.Errorf("creating item; no valid params count, params count: %d", len(params))
}

/*
Checks whether an item name length is zero
*/
func (db database) isItemValid(item string) (string, error) {
	if len(item) != 0 {
		return item, nil
	}

	return "", fmt.Errorf("checking the item \"%s\"; zero item name length", item)
}

/*
Checks whether the item is already placed
*/
func (db database) isGoodsPlaced(item string) bool {
	_, isPlaced := db[item]
	return isPlaced
}

/*
Checks whether a typed price is valid and returns the price and nil error if data is correct,
otherwise it return -1 and a corresponding error.
*/
func (db database) isPriceValid(priceParam string) (int, error) {
	parsedPrice, err := strconv.ParseInt(priceParam, 10, 32)
	if err != nil {
		return -1, fmt.Errorf("parsing a price argument: %s; %s", priceParam, err)
	}

	if parsedPrice <= 0 {
		return -1, fmt.Errorf("invalid price value: %d", parsedPrice)
	}

	return int(parsedPrice), nil
}
