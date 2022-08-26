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

start:
	air

test:
	go clean -testcache
	go test ./tests -v