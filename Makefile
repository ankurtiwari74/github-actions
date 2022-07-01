test: ## Run unittests
	@go test cmd/web/*.go

build: ## Build the binary file
	go build cmd/web/*.go

lambda: ## Build the zip archived binary to be uploaded as lambda
	GOOS=linux GOARCH=amd64 go build cmd/web/*.go
	zip -u main.zip main version

run: ## Run the app
	@go build -o Bookings cmd/web/*.go && ./Bookings