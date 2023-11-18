# technoStore
### Run Product Store Server

```bash
make run
```
Visit http://localhost:8080/swagger/index.html#/ for swagger doc of the project

### Run LRU
```bash
make lru
```
### Test

```bash
make test
```

# Architecture

In this project structure, the cmd directory contains the entry point of the application, which is the main.go file. The internal directory contains the implementation details of the application, which are divided into three layers: 
- api: The api layer contains the input/output logic of the application
- domain: The domain layer contains the business logic and entities of the application
- infrastructure: The infrastructure layer contains the implementation details

#Stack
I think pgx is amazing, I like working with it directly without database/sql middleman. It provides great features and performance.

# SQL Schema
### Migrations
run (if not installed)
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
To create a new migration
```
migrate create -ext sql -dir database/migrations -seq name_for_migration
```

and then edit the new files in `/database/migrations`.

You can also run the migration manually:

```
make migrate-up
```
for db down
```
make migrate-down
```

### Schema Dump (if use docker)
```
docker exec pg-docker pg_dump -U postgres --schema-only --verbose -s technoStore > ts_ddl.sql
```

# GoMock: for mocking definitions to test business layer
Install gomock
```
go install go.uber.org/mock/mockgen@latest
```
Generate mock datastore definition implementation
```
make mock-db
```
# API(Swagger doc generation)
```
make swag
```

# Deployment
### Docker build
```
docker build --progress=plain -f ./cmd/server/Dockerfile -t techno-store:latest -t techno-store:0.1 .
```

