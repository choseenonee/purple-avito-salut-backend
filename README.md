# go-gin-template
Template with 
<br>
Postgres (goose migrations)
<br>
Redis
<br>
Rabbit
<br>
Gin (swaggo swagger gen)
<br>
Docker, and docker-compose
<br>

launch redis and postgres:
<br>
``docker run --name postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres``
<br>
``docker run --name redis -d -p 6379:6379 redis``
<br>
``docker run -d --name jaeger -p 16686:16686 -p 14268:14268 jaegertracing/all-in-one:1.41``
<br>
goose migration: ``goose -dir deploy/migrations postgres "postgresql://postgres:postgres@localhost:5432/postgres" up``
<br>
launch main.go with cmd workdir, and [click](http://127.0.0.1:8080/swagger/index.html#/)
<br>




**_TRACING ONLY IN REFRESH/GETME handlers_**