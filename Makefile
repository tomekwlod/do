GOTEST=go test -cover
GOCOVER=go tool cover
.PHONY: tc
tc:
		$(GOTEST) -v -coverprofile=coverage.out ./...
		$(GOCOVER) -func=coverage.out
		$(GOCOVER) -html=coverage.out