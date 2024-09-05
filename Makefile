
docker-up:
	docker compose -f docker-compose.yaml up --build &

up-postgres:
	docker compose -f docker-compose.postgres.yaml up --build