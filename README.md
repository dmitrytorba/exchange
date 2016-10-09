#### setup
npm install

#### build
./build.sh

#### adding js libraries
npm install jquery --save

#### database
```
createdb exchange
psql -f schema.sql exchange
```

#### orderbook_db.go
This is where the orderbook is transcribed to the database. It synchronizes
// the in-memory book with database

#### orderbook.go
This is where the idea of an orderbook is put together. An
orderbook is an array of orders (buys or sells) and they must
be sorted and added/updated/deleted according to a set of rules.
The action of updating or deleting an order in the book is called
an execution, because the order is being partially or fully executed.

#### exchange.go
This is where the collcetion of orderbooks will reside and other exchange level
things will be handled.
