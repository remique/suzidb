# Suzi DB

SuziDB is a database system built for educational purposes. The main features are:

- SQL interface, including joins.
- Pluggable storage engines, including Bitcask written from scratch and in-memory.
- Iterator processing model, for query execution.

## Usage

```
> create table a(id int primary key, name text);
OK

> create table b(id int primary key, other text);
OK

> insert into a(id, name) values (1, 'alice');
OK

> insert into b(id, other) values (1, 'kowalski');
OK

> select * from a left join b on a.id = b.id;
+----+-------+----------+
| ID | NAME  | OTHER    |
+----+-------+----------+
|  1 | alice | kowalski |
+----+-------+----------+
```
