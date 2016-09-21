
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
  email varchar(254),
  token text,
  password text
);
grant all privileges on table users to exchange;
grant usage, select on sequence users_id_seq to exchange;

CREATE TYPE ordertype AS ENUM ('buy', 'sell');
create table if not exists orders
(
  id serial primary key,
  amount bigint,
  price bigint,
  order_type ordertype,
  username varchar(32)
);
grant all privileges on table orders to exchange;
grant usage, select on sequence orders_id_seq to exchange;