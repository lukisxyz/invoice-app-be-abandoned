#!make
include .env

DB_CONN = postgres://${DB_USER}:${DB_SECRET}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}

cmgr:
	bin/migrate create -ext sql -dir database/migration -seq ${name}

migup:
	bin/migrate -path database/migration -database "${DB_CONN}" -verbose up

migdown:
	bin/migrate -path database/migration -database "${DB_CONN}" -verbose down

migupx:
	bin/migrate -path database/migration -database "${DB_CONN}" -verbose up ${num}

migdownx:
	bin/migrate -path database/migration -database "${DB_CONN}" -verbose down ${num}

setupair:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

setupmigrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz -C bin/

setup: setupair setupmigrate

seed.product:
	./database/seeder/product.sh

run:
	bin/air
