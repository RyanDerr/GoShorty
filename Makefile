up: 
	docker compose up -d

down: 
	docker compose down
	docker rmi goshorty-api:local

test:
	go test -v ./... -cover

clean-db:
	docker exec -it redis redis-cli FLUSHALL
	docker exec -it user_db psql -U postgres -c "TRUNCATE users RESTART IDENTITY CASCADE;"

generate-spec:
	swag fmt -d ./
	swag init -g api/routes/router.go -o ./api/docs 

generate-cli:
	go build -o goshorty cmd/main.go