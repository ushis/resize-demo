SRC := $(shell find src -name '*.go')
BIN := resized

.PHONY: $(BIN)

$(BIN): $(SRC)
	gd -o $@

clean:
	rm -f $(BIN)
	gd clean

fmt:
	gd fmt -w2
