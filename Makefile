
start:
	@./scripts/dev.sh start
.PHONY: service

migrate:
	@./scripts/dev.sh migrate
.PHONY: migrate

up:
	@./scripts/dev.sh up
.PHONY: infras-up

down:
	@./scripts/dev.sh down
	
.PHONY: infras-down