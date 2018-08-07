NAME=$(notdir $(shell go list))
GOBIN:=$(shell (type vgo > /dev/null 2>&1) && echo vgo || echo go)

GO:=$(find . -name "*.go" -print)
LIBGO:=$(wildcard lib/*.go)
LIB:=$(LIBGO:.go=.so)

.SUFFIXES: .go .so
.go.so:
	$(GOBIN) build -buildmode=c-shared -o $@ $<

all: $(GO) $(LIB) test
	$(GOBIN) build

clean:
	rm -f $(NAME) $(wildcard lib/*.h) $(wildcard lib/*.so)

test:
	$(GOBIN) test

golint:
	$(GOBIN) list ./... | xargs -L1 golint

reviewdog:
	$(GOBIN) list ./... | xargs -L1 golint | reviewdog -f=golint -diff="git diff master"
