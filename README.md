# Event Store on PostgreSQL

The implementation of event sourced system that uses PostgreSQL as an event store. 

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
curl --location --request POST 'http://localhost:4012/place_order' \
--header 'Content-Type: application/json' \
--data-raw '{
    "customer_id": 123123,
    "price": 10.2
}'
```

Cancel order

```curl
curl --location 'http://localhost:4012/cancel_order' \
--header 'Content-Type: application/json' \
--data '{
    "order_id": "fe55e443-2426-437a-9656-f2daf01fa2f1"
}'
```