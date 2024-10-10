db:
	docker run --name postgres -p 8025:5432 -e POSTGRES_PASSWORD=dev -d  --network app postgres:15.6

d_build:
	docker build -t taskbackdev:1.0 .

d_run:
	docker rm -f taskbackdev && docker run --name taskbackdev --network app -p 8022:8081 -v $(pwd)/config-docker.yaml:/app/config.yaml:ro -d taskbackdev:local

