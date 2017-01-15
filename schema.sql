
do $$ begin
   if not exists
      (select * from pg_catalog.pg_user where usename = 'exchange') then
          create role exchange login password 'xNzoA3ZNfTe89Kqp2h';
          grant all on database exchange to exchange;   
       end if;
end $$;

create table if not exists users
(
  id serial primary key,
  username varchar(32),
  email varchar(254),
  email_token text,
  password text,
  btc bigint,
  ltc bigint
);
grant all privileges on table users to exchange;
grant usage, select on sequence users_id_seq to exchange;

--don't change the following index names or else the signup code will fail
CREATE UNIQUE INDEX unique_username on users (lower(username));
CREATE UNIQUE INDEX unique_email on users (lower(email));


CREATE TYPE ordertype AS ENUM ('buy', 'sell');
CREATE TYPE currency AS ENUM ('btc', 'ltc', 'usd');
create table if not exists orders
(
  id serial primary key,
  amount bigint,
  price bigint,
  order_type ordertype,
  username varchar(32),
  currency currency,
  timestamp timestamp
);
grant all privileges on table orders to exchange;
grant usage, select on sequence orders_id_seq to exchange;

create table if not exists executions
(
  id serial primary key,
  amount bigint,
  price bigint,
  order_type ordertype,
  filler varchar(32),
  username varchar(32),
  currency currency,
  timestamp timestamp
);
grant all privileges on table executions to exchange;
grant usage, select on sequence executions_id_seq to exchange;

CREATE TYPE exchange_name AS ENUM ('bitfinex');

create table if not exists bitfinex_trades_btcusd
(
  time_stamp timestamp primary key,
  time_recieved timestamp,
  price numeric,
  volume numeric
);
grant all privileges on table bitfinex_trades_btcusd to exchange;

-- ordercount == 0 means delete 
create table if not exists bitfinex_book_btcusd
(
  time_stamp timestamp primary key,
  price numeric,
  order_count bigint,
  order_type ordertype,
  volume numeric
);
grant all privileges on table bitfinex_book_btcusd to exchange;

create table if not exists gdax_book_btcusd
(
  time_stamp timestamp primary key,
  price numeric,
  order_type ordertype,
  volume numeric
);
grant all privileges on table gdax_book_btcusd to exchange;
  
create table if not exists gdax_trades_btcusd
(
  time_stamp timestamp primary key,
  time_recieved timestamp,
  price numeric,
  volume numeric
);
grant all privileges on table gdax_trades_btcusd to exchange;
