# GO SHOPPING CART with gRPC 

The Shopping Cart

In this gRPC micro service you can see basic features and patterns in Go:
* DDD Cart Aggregate
* Hexagonal Architecture
* Command Query Segregation CQS  
* SOLID
* Domain - Unit tests (TDD)
* Application test cases (BDD)
* Repository Pattern
* Transactional Persistence
* Middlewares

## Model

```

  ┌───────────────────┐           ┌───────────────────┐
  │       Cart        │<>─────────│     Product       │ 
  └───────────────────┘           └───────────────────┘ 

```
## How can I use it?

### Prerequisites
- Docker

### Installing

Download repository
```sh
git clone https://github.com/amelendres/go-shopping-cart
```


Build container
```sh
make build
```

## Run tests

Run the tests
```sh
make sh
make test
```


## Try it!
The shopping cart is running on `http://localhost:8050`

Check the services with [BloomRPC](https://github.com/uw-labs/bloomrpc) ,
importing the `proto/cart.proto`

  <img src="https://github.com/uw-labs/bloomrpc/raw/master/resources/editor-preview.gif" />




Compile *proto files* on Go
```sh
make sh
make gproto
```



### TODO

* Refactor Application services by ctx  
* Refactor CartServer with wire DI
* Refactor SRP CartServer Handlers
* Refactor product units to qty.  
* Add http request validation   
* Refactor Postgres Store to `SQLStore`
* Refactor `handler_test` case runner
* Optimize `updateCartTx` in order to just save the changes
* Add DB Migrations ... (goland migrate)
* Add Error response handler
* Add quality code checker
* Add domain events
* add event publisher
* add message publisher
    
    


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[MIT license](LICENSE.md)
