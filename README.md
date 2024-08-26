## Fulcrum project aims at simplifying deployment of applications at GRID.GG



### SQLITE3 acts as session storage
```
sqlite3 fulcrum.db
create table [sessions] (
token char(43) primary key,
data BLOB NOT NULL,
expiry TIMESTAMP(6) NOT NULL
);  
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```


