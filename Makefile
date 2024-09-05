.PHONY: run run-all run-client run-server run-migrate clean

run: run-all

run-all:
	$(MAKE) -j 2 run-migrate run-postgres run-server run_client 

run_client:
	npx nx dev Trio --port 3000

run-postgres:
	docker-compose -f ./compose.yml up -d

run-server:
	cd apps/server && \
	CompileDaemon -command="./trio"

run-migrate:
	cd apps/server && \
	go run migration/migrate.go

clean:
	docker stop trio-db && docker rm trio-db
	-pkill -f "CompileDaemon -command=./trio" # Stop the server
	-pkill -f "npx nx dev Trio" # Stop the web app