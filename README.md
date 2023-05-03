# Event Store on PostgreSQL

The implementation of event sourced system that uses PostgreSQL as an event store. 

## How to start

Start infras
```make
make up
```

Run migration(need to manual create database first)
```make
make migrate
```

Start service
```make
make start
```