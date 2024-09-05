.PHONY: run run-all run_client run_server clean

run: run-all

run-all:
	$(MAKE) -j 2 run-postgres run-server run_client 

run_client:
	npx nx dev Trio --port 3000

run-postgres:
	docker-compose -f ./compose.yml up -d

run-server:
	cd apps/server && \
	CompileDaemon -command="./trio"

clean:
	-pkill -f "CompileDaemon -command=./server" # Stop the server
	-pkill -f "npx nx dev Trio" # Stop the admin-portal app