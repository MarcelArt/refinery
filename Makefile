swag:
	@swag init --parseDependency --parseInternal

dev: swag
	@air serve

go: swag
	@go run main.go serve

migrate:
	@go run main.go migrate

migrate-force:
	@go run main.go migrate --drop

compose-down:
	@docker compose down

compose: compose-down
	@docker compose up -d

secret:
	openssl rand -base64 32

pull:
	@git pull

update: pull compose

license:
	@rm -rf THIRD_PARTY_LICENSES && go-licenses save ./... --save_path=./THIRD_PARTY_LICENSES