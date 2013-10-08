BIN := resized

export GOPATH := $(shell pwd)

.PHONY: $(BIN)

$(BIN): deps
	go build -v $@

deps:
	go get -d -v resizer/...

clean:
	rm -f $(BIN)

fmt:
	gofmt -l -w -tabs=false -tabwidth=2 src/resize{r,d}
