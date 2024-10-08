db:
	docker run --name postgres -p 8025:5432 -e POSTGRES_PASSWORD=dev -d  --network app postgres:15.6
