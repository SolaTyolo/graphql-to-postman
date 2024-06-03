export

update-postman-schema:
	go run cmd/update_collection_schema/main.go
.PHONY: update-postman-schema