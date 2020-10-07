build:
	docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml up -d
	go build -o echo-server
	./echo-server