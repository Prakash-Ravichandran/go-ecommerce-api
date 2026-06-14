# go-ecommerce-api server-server communication

<img width="962" height="537" alt="Image" src="https://github.com/user-attachments/assets/923cdcf2-79c5-479a-8c80-12993597a5a9" />

## API gateway :8080

[route REST /products service calls to gRPC /products
](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/c5fbb79f233f79d11751c70f17eee8dc2bd865f6)

## gRPC migration of /orders

- [orders proto rpc GetOrders, mock reponse
  ](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/897437ea9d5bc318caf0368f0bf0225403d53c86)

- [gRPC GetOrders Actual response &CreateOrders mockresponse
  ](https://github.com/Prakash-Ravichandran/go-ecommerce-api/commit/2cc8f82cd8a8d048b3dbf0dc85e2abc4f5429aee)

gRPC /CreateOrders JSON message input for postman

```json
{
  "customer_id": 49,
  "items": [
    {
      "product_id": 1,
      "quantity": 4
    },
    {
      "customer_id": 49,
      "product_id": 2,
      "quantity": 6
    }
  ]
}
```

response

```json
{
  "order": {
    "id": "21",
    "customer_id": "49",
    "created_at": "2026-06-14T13:27:41.273Z"
  }
}
```
