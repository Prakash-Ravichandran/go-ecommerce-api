# go-ecommerce-api

Use of internal folder in go [ref](https://www.bytesizego.com/blog/golang-internal-package)

- external users/usages of this package cannot have access to this internal folder.

## creating an http-server setup

- Add structs for http-server [commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/890847824557e7566fbfef66b78dcb09dd348170)
- Flow of request: user -> handler GET /products -> service getProducts -> DB SELECT \* FROM products
- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/7b9d8aab0df2ce17917d7367501be29e87dde9fe)

### Strcutured logging

- Instead of printing, we can use structured logging with levels - info, error amd with meta data
- Use global logger `slog.SetDefault(logger)` to replace `r.Use(middleware.Logger)`
  [structured logging](https://go.dev/blog/slog)

### Clean Layered Architecture

| |Transport| |
| |Service| |  
| |Repository||

Transport has a dependency of Service and Service depends on DB(repository)

- Transport layer contains http, grpc and it depends on both service and repository as a dependency
- Service layer contains the business logic and it depends on only repository as a dependency

<img width="1168" height="602" alt="Image" src="https://github.com/user-attachments/assets/103c11ef-9ef0-43d6-a422-d81b893a1d98" />

### when to use \* and & ?

Why return a pointer to the Handler ?

In NewHandler, you return \*handler. This is so your router (chi) can "own" an instance of the product logic.

By returning a pointer, you ensure that:

- Memory is efficient: **You aren't copying the whole handler struct every time**.
- State is shared: If the handler had a logger or a database connection, all routes using that handler would share the same one.

```go
func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}
```

[Golang pointers](https://youtu.be/2XEQsJLsLN0?si=bAQUEkC2mONMBqKk)

### add handler and service layers for /products route

[[products handler]](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/8124dc9ac98b430a8989e7486d82f39d347ca01b)

### custom json write package

[cutsom json write package](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/0ef7ebf458d1c180b5556add763fd28952ba4b56#diff-db5f40f068f65e740724c149d30ee31afe66fee767c4e663a8d864b5b7f5879a)

### Products service

[commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/7f6250ec98787c6ef7c04c6075b467bfaebd75dc)

- What is a service ?

The service doen't care about thr HTTP, JSON or headers. It only cares about the rules of the app (where a user can only list products if thet are active).

- what is a context ?

context package is Go's way of managing request lifecycles - cancellation, timeouts, tracing, tracing.

cancellation: If a user closes their browser while the server is still fetching products, the ctx will signal "Cancelled!" The service can then stop its work immediately, saving server resources.

- why two versions of ListProducts ?

handler version: It speaks HTTP. It extracts data from the http.Request, calls the service, and then formats the result into JSON for the http.ResponseWriter.
service version: It speaks Business Logic. It doesn't know what JSON is. It just knows how to go to the database (or wherever your products are) and return the data or an error.

#### Data Flow (request journey)

- The Router receives the HTTP request and sends it to productHandler.ListProducts.

-The Handler (The "Outer Shell") peels back the HTTP layer. It takes the r.Context() and passes it down.

- The Service (The "Brain") receives the context. It might check a database or a cache. It doesn't know about w or r (the response or request). It just returns data or an error.

- The Handler gets the result back. If there's an error, it decides: "This should be a 500 error." If it's successful, it encodes the data into JSON and sends it to the user.

#### Why use an Interface for the Service?

```go
type Service interface {
    ListProducts(ctx context.Context) error
}
```

This is a Dependency Injection pattern. By requiring an interface in your handler instead of the concrete svc struct, you gain two massive advantages:

Mocking for Tests: You can write a test for your Handler without needing a database. You just create a mockService struct that satisfies the Service interface and returns a fake error to see if your Handler handles it correctly.

Flexibility: If you decide to change your service logic (e.g., NewAdvancedService()), you don't have to change a single line of code in your handler.

#### The Power of Context (ctx)

In your service, you see ListProducts(ctx context.Context). Even though you aren't using ctx inside the function yet, it's vital.

Imagine ListProducts took 10 seconds to run because the database was slow. If the user gets tired and closes their browser after 2 seconds:

1. The HTTP server detects the disconnection.

2. The ctx is automatically "canceled."

3. Inside your service, you could check ctx.Err(). If it's canceled, you stop the database query mid-way.

Without Context, your server would keep working on a request that no one is listening to anymore!

### Database with sqlc

[commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/ddc8d6f77562725027ab08df4fa48b10c87ee588)

### Migrations

- use a package called goose to create a new sql migration [pkg.go.dev](https://pkg.go.dev/github.com/pressly/goose)

```go
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- History of changes to the database

- Database migrations are a way to incrementally modify your database schema. For example, adding new tables, altering existing tables, etc...You can write them as SQL files and execute them, or write them using libraries.
  [reddit ref](https://www.reddit.com/r/node/comments/90fo0t/whats_datadatabase_migration/#:~:text=Database%20migrations%20are%20a%20way,or%20write%20them%20using%20libraries.)

- goose up: what we are going to change in our database
- goose down: how we are going to rollback the change

- create a migrations sql file

```go
goose -s create create_products sql
```

#### what is dependency injection ?
