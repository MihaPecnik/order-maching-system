# Order Matching Engine (based on Stock Market principle):

## How to run application

  1. docker-compose build
  1. docker-compose up
  
  Or just:
  Run postgres cointainer from docker-compose.yml file.
  go build main.go [-parameters]
  
  parameters examples:
   * --migrate //migrate the databse
   * --postgres_url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" //databse connection
   * --populate //populate the database
   
## Testing

*  go test github.com/MihaPecnik/order-maching-system/handler -v
> 
> ok      github.com/MihaPecnik/order-maching-system/handler       0.395s

> PASS

> coverage: 100.0% of statements


## Api endpoints
* localhost:8080
###  Update Order's book
Given as input an id User, a ticker (eg. AAPL), a value, an int quantity and a
command (buy or sell) which is an order. Update an order’s book on the buy
and sell-side accordingly, return the quantities and prices that got matched as
a result of the order insertion.

If any error(e.g. electricity outage, server error) occurs, data will automatically rollback.


[Algorithm](https://www.youtube.com/watch?v=Kl4-VJ2K8Ik)
* Method:PUT 
* Path: /api/v1/api/v1/orderbook
```json
request:
{
    "user_id" :1,
    "value" : "199",
    "quantity": 7,
    "buy": false,
    "ticker": "APPL"
}
```
response:
```json
[
    {
        "value": "200.1",
        "quantity": 2
    },
    {
        "value": "200",
        "quantity": 2
    }
]
status:200
```

###  Get tickers bottom of the buy-side and the top of the sell-side
Given as input a ticker (return the bottom of the buy-side and the top of the
sell-side). Basically it just returns the "BUY" order with the highest value and it's quantity and order "SELL" with the lowest value and its quantity.


* Method:GET 
* Path: /api/v1/tickerinfo/{ticker}

response:
```json
{
    "buy": {
        "value": "20",
        "quantity": 15
    },
    "sell": {
        "value": "19",
        "quantity": 10
    }
}
status:200
```
