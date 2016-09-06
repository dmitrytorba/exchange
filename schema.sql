
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
  token varchar(1024),
  email varchar(254),
  password text,
  config text
);
grant all privileges on table users to exchange;

/* Theory of Operation 
# A buy order is made
# select all sells with price equal or below buy order price
# while buy order is not fill, fill the lowest priced orders until filled
# if all possible small orders exhausted, create new row in buy
# increase the relevant balance on user's row
*/
