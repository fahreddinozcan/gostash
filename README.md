To run PostgreSQL database, run 

```
docker run --name postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
```