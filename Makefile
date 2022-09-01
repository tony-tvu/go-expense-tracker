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
	cd web && npm run start	

watch: 
	cd web && npm run dev

air:
	air

start: 
	make -j 2 air watch

test:
	go clean -testcache
	go test ./smoketest -v
