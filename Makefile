BIN := cloverback

GOPATH := $(shell go env GOPATH)
GO_FILES := $(shell find . -name "*.go")

$(BIN): $(GO_FILES)
	gofumpt -w $(GO_FILES)
	go build -o $(BIN) cmd/main.go

install: $(GOPATH)/bin/$(BIN)
.PHONY: install

$(GOPATH)/bin/$(BIN): $(BIN)
	mv $(BIN) $(GOPATH)/bin/$(BIN)

clean:
	rm -f $(BIN)
.PHONY: clean
