# go-ecommerce-api

Use of internal folder in go [ref](https://www.bytesizego.com/blog/golang-internal-package)

- external users/usages of this package cannot have access to this internal folder.

## creating an http-server setup

- Add structs for http-server [commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/890847824557e7566fbfef66b78dcb09dd348170)
- Flow of request: user -> handler GET /products -> service getProducts -> DB SELECT \* FROM products
- [Commit](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/7b9d8aab0df2ce17917d7367501be29e87dde9fe)
