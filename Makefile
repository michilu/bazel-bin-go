GOBIN:=$(shell (type vgo > /dev/null 2>&1) && echo vgo || echo go)
PKG:=$(shell go list)
NAME:=$(notdir $(PKG))
PROJECT_SINCE:=$(shell git log --pretty=format:"%ad" --date=unix|tail -1)
AUTO_COUNT_SINCE:=$(shell echo $$(((`date -u +%s`-$(PROJECT_SINCE))/(24*60*60))))
AUTO_COUNT_LOG:=$(shell git log --since=midnight --oneline|wc -l|tr -d " ")

GO:=$(find . -name "*.go" -print)
LIBGO:=$(wildcard lib/*.go)
LIB:=$(LIBGO:.go=.so)

.SUFFIXES: .go .so
.go.so:
	$(GOBIN) build -buildmode=c-shared -o $@ $<

all: $(GO) $(LIB) test
	$(GOBIN) build -ldflags=" \
-X $(PKG)/meta.serial=$(AUTO_COUNT_SINCE).$(AUTO_COUNT_LOG) \
-X $(PKG)/meta.hash=$(shell git describe --always --dirty=+) \
-X \"$(PKG)/meta.build=$(shell LANG=en date -u +'%b %d %T %Y')\" \
"

clean:
	rm -f $(NAME) $(wildcard lib/*.h) $(wildcard lib/*.so)

test:
	$(GOBIN) test

golint:
	$(GOBIN) list ./... | xargs -L1 golint

reviewdog:
	$(GOBIN) list ./... | xargs -L1 golint | reviewdog -f=golint -diff="git diff master"
