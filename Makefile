DB_DIR      := database
COMPOSE     := docker compose -f $(DB_DIR)/docker-compose.yml
PRIMARY_DSN := postgres://postgres:postgres@localhost:5432/dvdrental?sslmode=disable
REPLICA_DSN := postgres://postgres:postgres@localhost:5433/dvdrental?sslmode=disable

.PHONY: up down start seed replica-init replica-init-1 replica-init-2 replica-init-3 \
        reset logs ps wait-primary wait-replica app-build app-run \
        frontend-install frontend-dev

start: ## Bring up containers, seed primary, initialise replica subscription
	$(MAKE) up
	$(MAKE) wait-primary
	$(MAKE) wait-replica
	$(MAKE) seed
	$(MAKE) replica-init

up: ## Start containers in the background
	$(COMPOSE) up -d

down: ## Stop containers
	$(COMPOSE) down

ps: ## Show container status
	$(COMPOSE) ps

logs: ## Tail all container logs
	$(COMPOSE) logs -f

wait-primary: ## Block until the primary accepts connections
	@echo "Waiting for primary..."
	@until $(COMPOSE) exec dvdrental-primary pg_isready -U postgres -d dvdrental -q 2>/dev/null; do \
		sleep 1; \
	done
	@echo "Primary is ready."

wait-replica: ## Block until the replica accepts connections
	@echo "Waiting for replica..."
	@until $(COMPOSE) exec dvdrental-replica pg_isready -U postgres -d dvdrental -q 2>/dev/null; do \
		sleep 1; \
	done
	@echo "Replica is ready."

seed: ## Load the dvdrental dataset into the primary
	psql "$(PRIMARY_DSN)" -f $(DB_DIR)/dvdrental.sql

replica-init: ## Set up logical replication (run steps 1-3 in order)
	$(MAKE) replica-init-1
	$(MAKE) replica-init-2
	$(MAKE) replica-init-3

replica-init-1: ## (Primary) Create replication user and publication
	psql "$(PRIMARY_DSN)" -f $(DB_DIR)/replication/01-replication-primary.sql

replica-init-2: ## (Replica) Apply schema
	psql "$(REPLICA_DSN)" -f $(DB_DIR)/replication/02-schema.sql

replica-init-3: ## (Replica) Create subscription
	psql "$(REPLICA_DSN)" -f $(DB_DIR)/replication/03-replication-replica.sql

app-build: ## Build the Go app
	cd app && go build -o ../bin/app .

app-run: app-build ## Run the Go app (reads app/config.yaml)
	./bin/app app/config.yaml

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run the Vite dev server (proxies /api -> :8080)
	cd frontend && npm run dev

reset: ## Tear down containers and wipe all data directories
	$(COMPOSE) down
	rm -rf $(DB_DIR)/datadir/dvd-pgdata $(DB_DIR)/datadir/dvd-replica-pgdata
