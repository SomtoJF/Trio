.PHONY: run run-all run-client run-server run-db-migrate swagger-migrate clean

run: run-all

run-all:
	$(MAKE) -j 2 run-postgres run-db-migrate swagger-migrate run-server run_client 

run_client:
	npx nx dev Trio --port 3000

run-postgres:
	docker-compose -f ./compose.yml up -d

run-server:
	cd apps/server && \
	CompileDaemon -command="./trio"

migrations:
	$(MAKE) run-db-migrate && $(MAKE) swagger-migrate

run-db-migrate:
	cd apps/server && \
	go run migration/migrate.go

swagger-migrate:
	cd apps/server && \
	swag init --parseDependency true

clean:
	docker stop trio-db && docker rm trio-db
	docker stop trio-qdrant && docker rm trio-qdrant
	-pkill -f "CompileDaemon -command=./trio" # Stop the server
	-pkill -f "npx nx dev Trio" # Stop the web app