
include .env

confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# ============================================================================= #
# BUILD
# ============================================================================= #
build:
	go build -o travex ./cmd/*.go

start:
	go run ./cmd/*.go

# ============================================================================= #
# QUALITY CONTROL
# ============================================================================= #
audit: vendor
	@echo 'Tidying and verifying module dependencies...'
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


migrateup: 
	 @echo 'Running up migrations...'
	migrate -path db/migrations -database ${POSTGRES_URL} -verbose up 

seed: 
	 @echo 'Running up seeds..'
	 migrate -path db/seed -database ${POSTGRES_URL} -verbose up 
	 go run ./db/seed/*.go

migratedown: 
	 migrate -path db/migrations -database {POSTGRES_URL} -verbose down


.PHONY: migrateup migratedown audit build start seed