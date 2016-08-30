
do $$ begin
   if not exists
      (select * from pg_catalog.pg_user where username = 'exchange') then
          create role exchange login password 'xNzoA3ZNfTe89Kqp2h';
          grant all on database exchange to exchange;   
       end if;
end $$;

create table if not exists users
(
  id serial primary key,
  token varchar(1024),
  email varchar(254),
  password text,
  config text
);
grant all privileges on table users to exchange;

CREATE TYPE order_type AS ENUM ('sell', 'buy');
CREATE TYPE currency AS ENUM ('ltc');

CREATE TABLE orders (
	order_type order_type,
	currency currency,
	amount bigint,
	price bigint,
	user_id int,
	age timestamp
);
grant all privileges on table orders to exchange;

/* sells should be sorted with the smallest price at the top, and if two match in price, oldest at the top */
CREATE INDEX on orders(currency, price DESC, age ASC) WHERE order_type = 'sell';

/* buys should be sorted with the largest price at the top, and if two match in price, oldest at the top */
CREATE INDEX on orders(currency, price ASC, age ASC) WHERE order_type = 'buy';

/* Theory of Operation 
# A buy order is made
# select all sells with price equal or below buy order price
# while buy order is not fill, fill the lowest priced orders until filled
# if all possible small orders exhausted, create new row in buy
# increase the relevant balance on user's row
*/
