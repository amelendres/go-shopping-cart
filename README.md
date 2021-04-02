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

* End 2 end test scenarios
  * Refactor Postgres Store to `SQLStore`
  * Add test ENV with sqlite
  * Refactor `handler_test` case runner
* Optimize `updateCartTx` in order to just save the changes
* Add application context
* Use internal packages
* Add DB Migrations
* Add request validator
* Add Error handler
* Add quality code checker
* Add domain events
  * add event dispatcher
    
    


## Authors

* **Alfredo Melendres** -  alfredo.melendres@gmail.com

[MIT license](LICENSE.md)
