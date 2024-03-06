

typo:
	@typos

setup-pre-commit:
	@pre-commit install

test-v-http:
	@go test ./http -v 
