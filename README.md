# GO SHOPPING CART gRPC 

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

## How can I use it?

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

Generate Proto
```
make sh
make gproto
make test
```

### HOW TO RUN

### Prerequisites
- BloomRPC



## TODO

* Add end 2 end test scenarios
  * Refactor Postgres Store to `SQLStore`
  * Add test ENV with sqlite
  * Refactor `handler_test` case runner
* Refactor cart `Products` as `ProductLines` 
  * Update database `product` to `product_lines`
  * Add auto incremental `id`
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
