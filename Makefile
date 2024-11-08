up: 
	docker compose up -d

down: 
	docker compose down
	docker rmi goshorty-api:local

clean-db:
	docker exec -it db redis-cli FLUSHALL

generate-spec:
	swag init -g api/main.go -o ./api/docs

generate-cli:
	go build -o goshorty cmd/main.go