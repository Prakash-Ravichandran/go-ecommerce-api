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

### add handler and service layers for /products route

[[products handler]](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/8124dc9ac98b430a8989e7486d82f39d347ca01b)
