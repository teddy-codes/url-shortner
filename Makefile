.PHONY: test up postgres

test:
	modd

# Postgres is required for the application which is why it's a pre-rec to `up`.
# Eventually, I will put this in a docker-compose file, but I am lazy
postgres:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=url_shortner postgres

up: postgres
	go run ./cmd/api