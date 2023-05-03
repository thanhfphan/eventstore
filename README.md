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