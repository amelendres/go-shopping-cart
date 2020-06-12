# GO SHOPPING CART API

### Prerequisites
- Docker

### Installing

Download repository
```
git clone https://github.com/amelendres/go-shopping-cart
```


Build container
```
make build
```

### HOW TO RUN 

**ENDPOINTS**

* Add a product to Cart
```
curl --location --request POST 'http://localhost:8050/carts/c50bf7b3-95d5-48fa-8b0d-691e3f40c1f9/products' \
--header 'Content-Type: application/json' \
--data-raw '{
  "id": "4e45b227-6a79-44ee-8cf0-da21508a4f8a",
  "name": "Dress",
  "price": 299.50,
  "units": 1
}'
```

* Get products from cart
```
curl --location --request GET 'http://localhost:8050/carts/c50bf7b3-95d5-48fa-8b0d-691e3f40c1f9/products' \
--header 'Content-Type: application/json' \
```


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[MIT license](LICENSE.md)
