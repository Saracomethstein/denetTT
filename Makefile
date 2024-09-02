SERVER_NAME=build/server
CLIENT_NAME=build/client

.PHONY: build clean server client deps

deps:
	@echo "==> Installing dependencies..."
	go mod tidy

build: deps
	@echo "==> Building the application..."
	mkdir build
	go build -o $(SERVER_NAME) server/cmd/server/main.go
	go build -o $(CLIENT_NAME) client/cmd/client/main.go

server: 
	@echo "==> Running the server..."
	./$(SERVER_NAME)

client:
	@echo "==> Running the client..."
	./$(CLIENT_NAME)

clean:
	@echo "==> Cleaning up..."
	go clean
	rm -f $(SERVER_NAME)
	rm -f $(CLIENT_NAME)
	rm -rf build
	