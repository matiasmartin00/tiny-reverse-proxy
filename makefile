# General configuration
IMAGE=hashicorp/http-echo

# Number of servers to run
NUM_USERS?=2
NUM_ORDERS?=2
NUM_PRODUCTS?=4

# Prefixes for each server
USER_PREFIX=user_server_
ORDER_PREFIX=order_server_
PRODUCTS_PREFIX=products_server_

# Base ports for each server, each server will be exposed in a port starting from this base port
PORT_BASE_USERS=5000
PORT_BASE_ORDERS=6000
PORT_BASE_PRODUCTS=7000

# Start servers
run:
	@echo "🚀 Init servers..."
	$(MAKE) start_servers TYPE=users PREFIX=$(USER_PREFIX) NUM=$(NUM_USERS) PORT_BASE=$(PORT_BASE_USERS)
	$(MAKE) start_servers TYPE=orders PREFIX=$(ORDER_PREFIX) NUM=$(NUM_ORDERS) PORT_BASE=$(PORT_BASE_ORDERS)
	$(MAKE) start_servers TYPE=products PREFIX=$(PRODUCTS_PREFIX) NUM=$(NUM_PRODUCTS) PORT_BASE=$(PORT_BASE_PRODUCTS)

start_servers:
	@for i in $(shell seq 1 $(NUM)); do \
		PORT=$$(($(PORT_BASE) + $$i)); \
		NAME=$(PREFIX)$$i; \
		echo "▶️ Starting $$NAME on port $$PORT..."; \
		docker run -p $$PORT:5678 -d --name $$NAME $(IMAGE) -text="Hello from $$NAME"; \
	done

# Stop servers
stop:
	@echo "🛑 Stoping servers..."
	$(MAKE) stop_servers PREFIX=$(USER_PREFIX) NUM=$(NUM_USERS)
	$(MAKE) stop_servers PREFIX=$(ORDER_PREFIX) NUM=$(NUM_ORDERS)
	$(MAKE) stop_servers PREFIX=$(PRODUCTS_PREFIX) NUM=$(NUM_PRODUCTS)

stop_servers:
	@for i in $(shell seq 1 $(NUM)); do \
		NAME=$(PREFIX)$$i; \
		echo "⏹️ Shutting $$NAME..."; \
		docker stop $$NAME 2>/dev/null || true; \
		docker rm $$NAME 2>/dev/null || true; \
	done

# Restart servers
restart: stop run

# Show status of servers
status:
	@docker ps --format "table {{.Names}}\t{{.Status}}"

# Clean all servers
clean: stop
