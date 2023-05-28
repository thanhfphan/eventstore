# Event Store on PostgreSQL

The implementation of event sourced system that uses PostgreSQL as an event store. This repository should use as a template.

## How to start

Start infras

```bash
./scripts/dev.sh up
```

Run migration(need to manual create database first)

```bash
./scripts/dev.sh migrate
```

Start service

```bash
./scripts/dev.sh start
```

Other commands

```bash
./scripts/dev.sh help
```

## Testing

Place order

```curl
curl --location --request POST 'http://localhost:4012/order/place' \
--header 'Content-Type: application/json' \
--data-raw '{
    "customer_id": 123123,
    "price": 10.2
}'
```

Cancel order

```curl
curl --location --request POST 'http://localhost:4012/order/cancel' \
--header 'Content-Type: application/json' \
--data '{
    "order_id": "fe55e443-2426-437a-9656-f2daf01fa2f1"
}'
```

Get order

```curl
curl --location --request GET 'http://localhost:4012/order/2d8cfbdc-fc57-4136-90b4-ba3211ed9dfb' \
--header 'Content-Type: application/json'
```