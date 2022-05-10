reload-server:
	gin --appPort 3000 --port 5000 --immediate

compose-down:
	docker-compose down
compose-up:
	docker-compose up -d

compose-build:
	docker-compose up -d --build

.PHONY: reload-server,
		compose-down,
		compose-up,
		compose-build
