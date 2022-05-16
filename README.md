# Recipes API simple REST API example
A REST API examples for simple todo application with Go

It is a just simple tutorial or example for making simple REST API with Go using **Gin Web Framework** and **mongo-driver** (An ORM of the mongo for Go)

## Installation & Run
```bash
# Download this project
git clone https://github.com/Andreis3/recipes-api.git
go mod download
```
## Frameworks
- [Mongo-driver](https://github.com/mongodb/mongo-go-driver)
- [Gin Web Framework](https://gin-gonic.com/)
- [go-redis](https://github.com/go-redis/redis)
- [session](https://github.com/gin-contrib/sessions)
- [session-redis](https://github.com/gin-contrib/sessions#redis)
- [go-swagger](https://goswagger.io/install.html)
- [postman](https://www.postman.com/)

## Make command
```bash
make reload-server # run reload server 
make compose-up # up mongo, redis, prometheus and grafana
make compose-down # down mongo, redis, prometheus and grafana
make compose-build # build mongo, redis, prometheus and grafana
make swagger-up # up swagger
```
## Install swagger
 - https://goswagger.io/install.html
 - swagger version

## API

#### /recipes
* `GET` : Get all recipes
* `POST` : Create a new recipe

#### /recipes/:id
* `PUT` : Update a recipe
* `DELETE` : Delete a recipe
* `GET` : Get a recipe

#### /recipes/search?tag=:tags
* `GET` : Get all recipes by tags

#### /signin
* `POST` : User signed in

#### /signout
* `POST` : User signed out


## Prometheus & Grafana

- http_requests_total
- http_request_duration_seconds
- http_status_codes_total

```
# Prometheus
- configure prometheus data source http://localhost:9090
```

Tutorial dashboard
- [Tutorial](https://www.programmingwithwolfgang.com/create-grafana-dashboards-with-prometheus-metrics/#:~:text=To%20create%20your%20own%20Grafana,and%20then%20select%20Add%20Query.&text=Select%20Prometheus%20as%20your%20data,over%20the%20last%2010%20minutes.)

*create app.env*
```
PORT=8080
MONGO_URI=mongodb://root:root@localhost:27017/test?authSource=admin
MONGO_DB=demo
COLLECTION_RECIPES=recipes
COLLECTION_USERS=users
REDIS_URI=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
X_API_KEY=1234560abcdetghig
```
