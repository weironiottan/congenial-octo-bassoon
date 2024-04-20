## API Reference

#### Get all orders from the Order Up Service

```http
  GET /orders
```

HTTP 200 OK Response:
```json
[
  {
  "id": "order-1234",
  "customerEmail": "martingarrix@email.com",
  "lineItems": [
    {
      "description": "Item 1",
      "priceCents": 100,
      "quantity": 1
    }
  ],
  "status": "2"
  }
]
```


#### Get all orders from the Order Up Service filtered by the orders' status
```http
  GET /orders/orders?status={orderStatus}
```

| Parameter     | Type     | Description                 |
| :------------ | :------- | :-------------------------- |
| `orderStatus` | `string` | pending, charged, fulfilled |

HTTP 200 OK Response:
```json
[
  {
  "id": "order-1234",
  "customerEmail": "martingarrix@email.com",
  "lineItems": [
    {
      "description": "Item 1",
      "priceCents": 100,
      "quantity": 1
    }
  ],
  "status": "2"
  }
]
```


#### Get an order by its order id
```http
  GET /orders/${id}
```

| Parameter | Type     | Description                                                   |
| :-------- | :------- | :------------------------------------------------------------ |
| `id`      | `string` | **Required** Must match an order id format such as: order-123 |

HTTP 200 OK Response:
```json
{
  "id": "order-1234",
  "customerEmail": "martingarrix@email.com",
  "lineItems": [
    {
      "description": "Item 1",
      "priceCents": 100,
      "quantity": 1
    }
  ],
  "status": "2"
}

```


#### Post a order

```http
  POST /orders
```

Post Order Body:
```json
{
  "customerEmail": "martingarrix@email.com",
  "lineItems": [
    {
      "description": "Item 1",
      "priceCents": 100,
      "quantity": 1
    }
  ]
}
```

HTTP 201 Created Response:

```json
{
  "id": "order-1234",
  "customerEmail": "martingarrix@email.com",
  "lineItems": [
    {
      "description": "Item 1",
      "priceCents": 100,
      "quantity": 1
    }
  ],
  "status": "2"
}

```

#### Charge the Order. Note that all fields in the post body are required
```http
  POST /orders/${id}/charge
```

| Parameter | Type     | Description                                                   |
| :-------- | :------- | :------------------------------------------------------------ |
| `id`      | `string` | **Required** Must match an order id format such as: order-123 |

Post Order Body:
```json
{
  "cardToken": "tokenized-credit-card-number"
}
```

HTTP 200 OK Response:
```json
{
  "chargedCents": 100,
}

```

#### Refund the Order. Note that all fields in the post body are required
##### Will only refund amount if the order status is not fulfilled
```http
  POST /orders/${id}/cancel
```

| Parameter | Type     | Description                                                   |
| :-------- | :------- | :------------------------------------------------------------ |
| `id`      | `string` | **Required** Must match an order id format such as: order-123 |

Post Order Body:
```json
{
  "cardToken": "tokenized-credit-card-number"
}
```

HTTP 200 OK Response:
```json
{
  "refundAmount": 100,
  "id": "order-1234"
}

```

#### Fulfill will change the order status to fulfilled
##### Will only change the status if the order is charged
##### This is an idempotent call, meaning you can call this endpoint on an already fulfilled order
##### with no change in status if its already status fulfilled

```http
  PUT /fulfill
```

| Parameter | Type     | Description                                                   |
| :-------- | :------- | :------------------------------------------------------------ |
| `id`      | `string` | **Required** Must match an order id format such as: order-123 |

Post Order Body:
```json
{
  "description": "item 1",
  "quantity": 2,
  "id": "order-1234"
}
```

HTTP 200 OK Response:
```json
{
  "id": "order-1234"
  "status": "2"
}
```

### All Possible error codes that can be expected from this service

HTTP 500 Internal Server Error:
```
Description: This is a general error meaning the request was unable to be made, this could be a network error or a failure to provide the expected request
```

HTTP 409 Conflict:
```
Description: This error occurs when there is a conflict in a resource. Most of the time it happens when the server is unable to change the order status
```

HTTP 404 Not Found:
```
Description: The resource requested was not found. This generally the order(s) requested are not present in the database
```

HTTP 400 Bad Request:
```
Description: The request that the service got was not the request it expected. This is a validation error. Please check that the request you are sending is correct before attempting again
```

