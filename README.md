## Setup

**[1]** Clone this repo and let refer path of `secret-server` as PROJECT_ROOT.
```
git clone https://github.com/nurali/secret-server.git
```

**[2]** Build secret-service docker image as below. This will do `go build` secret-service and then will do `docker build` to create docker image **secret-service:latest**.
```
cd $PROJECT_ROOT/secret-service
make build-ci
```

**[3]** Run secret-server with docker-compose as below.  This will start two docker container with name `secret-server_app_1` and `secret-server_db_1`.
```
cd $PROJECT_ROOT
docker-compose up -d
```

**[4]** Check secret-service and postgresql container are running fine.
```
docker ps

CONTAINER ID        IMAGE                                           COMMAND                  CREATED             STATUS              PORTS                    NAMES
50ef275e206b        secret-service:latest                           "./secret-service"       14 minutes ago      Up 14 minutes                                secret-server_app_1
1941ed5e2497        registry.centos.org/postgresql/postgresql:9.6   "container-entrypo..."   14 minutes ago      Up 14 minutes       0.0.0.0:5432->5432/tcp   secret-server_db_1
```


## Try

**[1]** Create Secret, with expireAfterViews=10
```
curl -X POST http://localhost:8080/api/secret --data '{"secret":"hello secret", "expireAfterViews":5, "expireAfter":5}'

{"hash":"7146a753-c9d9-481b-9a3c-ada179cb32b0","secretText":"hello secret","createdAt":"2019-07-31T12:19:24Z","expiresAt":"2019-07-31T12:24:24Z","remainingViews":5}
```

**[2]** Get Secret, for above hash, note remainingViews=4 
```
curl http://localhost:8080/api/secret/7146a753-c9d9-481b-9a3c-ada179cb32b0

{"hash":"7146a753-c9d9-481b-9a3c-ada179cb32b0","secretText":"hello secret","createdAt":"2019-07-31T12:19:24Z","expiresAt":"2019-07-31T12:24:24Z","remainingViews":4}
```
You can make `GET /api/secret` multiple times for same `hash`.  Once `remainingViews=0`, after that more call to Get endpoint will return error.

**[3]** Open broswer and see **Prometheus metrics** at http://localhost:8080/metrics search metrics with prefix **secret_service_**.


## Notes
- We can enhance to support multiple MIME types.  Currently `application/json` is supported.

- Currently, `expireAfter` value is NOT used (/considered)  for `GET /api/secret` endpoint.
