# buycoffee-order

## SQL

```bash
psql --username=postgres
```


```bash
psql -h 127.0.0.1 -p 32768 --username=buycoffeeuser -d buycoffee
```

```sql
CREATE DATABASE buycoffee;
CREATE USER buycoffeeuser;
ALTER USER buycoffeeuser with encrypted password 'supersecretpass';
GRANT ALL PRIVILEGES on database buycoffee to buycoffeeuser;

\c buycoffee buycoffeeuser

CREATE TABLE "order"(
   order_id serial PRIMARY KEY,
   user_id integer NOT NULL,
   item_name VARCHAR (355) NOT NULL,
   created_on TIMESTAMP NOT NULL
);

INSERT INTO "order"(user_id, item_name, created_on*)
VALUES
   (1, 'cappucino', NOW());
```