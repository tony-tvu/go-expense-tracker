install:
	go mod download
	cd web && npm install

build_prod:
	cd web && npm install
	cd web && npm run build
	rm -r web/node_modules
	go build -o ./bin/main main.go

build:
	cd web && npm run build
	go build -o ./bin/main main.go

api:
	go run main.go

client:
	cd web && npm start

start:
	air

test:
	go clean -testcache
	go test ./tests/unit -v
	go test ./tests/integration -v
