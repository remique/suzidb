# Suzi DB

```
> create table mytable(id int primary key, name text);
OK

> insert into mytable(id, name) values (1, 'alice');
OK

> insert into mytable(id, name) values(2, 'bob');
OK

> select * from mytable;
+----+-------+
| ID | NAME  |
+----+-------+
|  1 | alice |
|  2 | bob   |
+----+-------+
```
