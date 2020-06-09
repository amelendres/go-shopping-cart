# GO SHOPPING CART API

### Prerequisites
- Docker

### Installing

Download repository
```
git clone [repository]
```


Build

The first installation 
```
make build
```

### HOW TO RUN 

Run Tests: 
```
make test
```

Server Start: 
```
make start
```

Examples:

Add product

```
curl --location --request GET 'http://localhost:5000/products'
```

Get products
```
curl --location --request POST 'http://localhost:5000/products' \
--header 'Content-Type: application/json' \
--data-raw '{
  "Name": "Pants",
  "Qty": 1
}'

```


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[MIT license](LICENSE.md)
