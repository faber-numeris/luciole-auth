generate:
	@echo "Generating files..."
	go run github.com/ogen-go/ogen/cmd/ogen@latest --target api/gen --clean api/openapi.yaml