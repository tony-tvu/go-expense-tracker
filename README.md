# Go Expense Tracker

Expense tracker application written in Go, React, and MongoDB. It uses [Teller.io](https://teller.io/) api to automatically retrieve user financial transactions, however this is optional and can be used as a standalone expense tracking app where users can manually add transactions.


## Prerequisites
- Go
- Node
- Docker Desktop


## 1. Clone the repository


```bash
git clone https://github.com/tony-tvu/go-expense-tracker.git
cd go-expense-tracker
```

## 2. Set up environment variables and certificates

```bash
cp .env.example .env
```

Copy `.env.example` to a new file called `.env`. If you plan to use this app with Teller.io, be sure to fill in your teller application ID for `REACT_APP_TELLER_APPLICATION_ID` and place your certificates in the /certificate directory. Teller should've provided a `certificate.pem` and `private_key.pem` when you created an account and provided your app information. 

## 3. Start docker

```bash
docker compose up
```

## 3. Run App
```bash
make start
```
<br/>

### Special instructions for Windows
Install node dependencies
```bash
npm install
```
Build frontend
```bash
npm run build
```
Run App
```bash
go run main.go
```