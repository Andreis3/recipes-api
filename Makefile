reload-server:
	gin --appPort 8080 --port 5000 --immediate

compose-down:
	docker-compose -f docker-compose.yaml down
compose-up:
	docker-compose -f docker-compose.yaml up -d

compose-build:
	docker-compose -f docker-compose.yaml up -d --build

swagger-run:
	swagger serve ./swagger.json

.PHONY: reload-server,
		compose-down,
		compose-up,
		compose-build
		swagger-run
