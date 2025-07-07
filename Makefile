BINARY_NAME=terminal-cli
DB_NAME=notes_db.db

# Objetivo principal
.PHONY: all build clean run deps

all: build

# Compilar el binario
build:
	go build -o $(BYNARY_NAME) . -ldflags="-s -w"

# Limpiar archivos generados
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(DB_NAME)

# Ejecutar el binario
run: build
	./$(BINARY_NAME)

# Instalar y actualizar dependencias
deps:
	go mod tidy
	go get -u all
