.PHONY: generate-swagger-docs

generate-swagger-docs:
	swag init --generalInfo=./internal/app/app.go --parseInternal --parseDependency --output=./internal/docs/v1