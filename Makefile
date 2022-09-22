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

# requires gin installed globally: `go install github.com/codegangsta/gin@latest`
api:
	gin --port 8080 run main.go

client:
	cd web && npm run start	

watch: 
	cd web && npm run dev

start: 
	make -j 2 api watch

test:
	go clean -testcache
	go test ./tests -v
