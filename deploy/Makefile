# Define the docker-compose command for convenience
DC = docker-compose
DOC = docker
GIT = git
KUBE = kubectl

# Define the service names as variables for easier reference
APP_SERVICE = app
POSTGRES_SERVICE = postgres-database
REDIS_SERVICE = redis-database

# Default target executed when no arguments are given to make.
default: start

# Target to start the application
start:
	$(DC) --env-file ../.env up

# Target to perform a soft restart (without removing source images)
soft-restart:
	$(DC) stop
	$(DC) --env-file ../.env up

# Target to perform a hard restart (with removal of source images)
hard-restart: update
	$(DC) down
	$(DOC) rmi deploy-app
	$(DC) --env-file ../.env up --build

# Helper target to stop all services
stop:
	$(DC) stop

# Helper target to remove all containers, networks, and images created by up.
clean:
	$(DC) down --rmi all

# Helper target to view logs
logs:
	$(DC) logs -f

update:
	$(GIT) pull

kuber:
	$(KUBE) apply -f admin-config-map.yaml
	$(KUBE) apply -f master-postgres-pv.yaml
	$(KUBE) apply -f master-postgres-deployment.yaml
	$(KUBE) apply -f master-postgres-service.yaml
	$(KUBE) apply -f trace-deployment.yaml
	$(KUBE) apply -f trace-service.yaml
	$(KUBE) apply -f admin-panel-deployment.yaml
	$(KUBE) apply -f admin-panel-service.yaml

kuber-clean:
	$(KUBE) delete all --all

.PHONY: default start soft-restart hard-restart stop clean logs update kuber