
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# ============================================================================= #
# BUILD
# ============================================================================= #
build:
	go build -o pera ./cmd/*.go

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
	 migrate -path db/migrations -database "postgresql://postgres:password@localhost:15432/travex?sslmode=disable" -verbose up 

migratedown: 
	 migrate -path db/migrations -database "postgresql://postgres:password@localhost:15432/travex?sslmode=disable" -verbose down


.PHONY: migrateup migratedown audit build