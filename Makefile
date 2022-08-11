install:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
	npm install

build:
	npm run build
	GOARCH=amd64 GOOS=darwin go build -o goexpense main.go

start:
	./start.sh