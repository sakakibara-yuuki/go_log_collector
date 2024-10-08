# Go でsample.logからlogを収集する
## How to run

docker composeでpostgresを起動する
```bash
docker compose up -d
```
起動したら,
```bash
go run main.go sample.log
```
でlogを収集することができる

```bash
dev ?1 ❯ docker compose exec db psql -U root -d logdb
psql (17.0 (Debian 17.0-1.pgdg120+1))
Type "help" for help.

logdb=# \dt
       List of relations
 Schema | Name  | Type  | Owner
--------+-------+-------+-------
 public | users | table | root
(1 row)

logdb=# select * from users;
 id | age |  name   |      role
----+-----+---------+-----------------
  1 |  22 | tarou   | student
  2 |  23 | zirou   | student
  3 |  24 | saburou | student
  4 |  25 | mike    | mentor
  5 |  26 | make    | mentor
(5 rows)

logdb=#
```
Got it!
