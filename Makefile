update-mocks:
	@rm -rf mocks
	@mkdir mocks
	@touch mocks/.coverignore
	@go run cmd/mockery/main.go --quiet

mockscfg:
	go run cmd/mockery/main.go showconfig