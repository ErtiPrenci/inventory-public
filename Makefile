# Nombre del binario final
APP_NAME=inventory-backend

# ==============================================================================
# Development commands
# ==============================================================================

## run: Runs the application (Assumes it loads .env internally or you have it in the environment)
run:
	@echo "Initializing server..."
	go run cmd/api/main.go

## build: Compiles the binary for production (Linux/Server)
build:
	@echo "Compiling..."
	go build -o bin/$(APP_NAME) cmd/api/main.go

## clean: Removes binaries and temporary files
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin/

## test: Runs the tests
test:
	@echo "Running tests..."
	go test -v ./...

## tidy: Downloads libraries and cleans go.mod
tidy:
	@echo "Organizing dependencies..."
	go mod tidy

build-lambda:
	@echo "Limpiando compilaciones previas..."
	@rm -rf build
	@mkdir -p build
	@echo "Compilando binario bootstrap (ARM64)..."
	# CGO_ENABLED=0 es vital para binarios que corren en Lambda
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o build/bootstrap cmd/api/main.go
	@echo "Dando permisos de ejecución..."
	@chmod +x build/bootstrap
	@echo "Creando paquete ZIP..."
	# El flag -j (junk paths) mete el binario en la raíz del zip
	@cd build && zip -j main.zip bootstrap
	@echo "Listo: build/main.zip creado correctamente."
# ==============================================================================
# Help
# ==============================================================================

## help: Shows available commands
help:
	@echo "Usage: make [command]"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: run build clean test tidy help