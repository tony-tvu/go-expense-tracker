test:
	go clean -testcache
	go test ./tests/... -v
	
build: 
	npm run build
	go build -o ./bin/main main.go

node_modules: package.json Makefile
	npm install
	touch $@

start: node_modules build
	bin/main
