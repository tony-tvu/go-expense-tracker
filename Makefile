install:
	go mod download
	cd web && npm install

build:
	cd web && npm run build
	go build -o ./bin/main cmd/main.go

run_binary:
	./bin/main

api:
	air

client:
	cd web && npm start
	
start:
	make -j2 api client