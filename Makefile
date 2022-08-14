install:
	go mod download
	cd web && npm install

install_pw:
	go get github.com/playwright-community/playwright-go
	go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps

build:
	cd web && npm run build
	go build -o ./bin/main main.go

build_clean:
	cd web && npm install
	cd web && npm run build
	rm -r web/node_modules

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
	go test ./...