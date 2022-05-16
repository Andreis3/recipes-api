reload-server:
	gin --appPort 8080 --port 5000 --immediate

compose-down:
	docker-compose down
compose-up:
	docker-compose -f docker-compose.yaml up -d

compose-build:
	docker-compose up -d --build

.PHONY: reload-server,
		compose-down,
		compose-up,
		compose-build
