install:
	go mod download
	cd web && npm install

build:
	cd web && npm install
	cd web && npm run build
	rm -r web/node_modules
	go build -o ./bin/main main.go

build_local:
	cd web && npm run build
	go build -o ./bin/main main.go

run_binary:
	./bin/main

api:
	air

client:
	cd web && npm start
	
start:
	make -j2 api client

test:
	go clean -testcache
	go test ./tests -v