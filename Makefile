PROJECT_SINCE:=$(shell git log --pretty=format:"%ad" --date=unix|tail -1)
AUTO_COUNT_SINCE:=$(shell echo $$(((`date -u +%s`-$(PROJECT_SINCE))/(24*60*60))))
AUTO_COUNT_LOG:=$(shell git log --since=midnight --oneline|wc -l|tr -d " ")
COMMIT:=4b825dc
REVIEWDOG:=| reviewdog -efm='%f:%l:%c: %m' -diff="git diff $(COMMIT) HEAD"

GOBIN:=$(shell (type vgo > /dev/null 2>&1) && echo vgo || echo go)
PKG:=$(shell $(GOBIN) list)
NAME:=$(notdir $(PKG))
GOLIST:=$(shell $(GOBIN) list ./...)
GODIR:=$(patsubst $(PKG)/%,%,$(wordlist 2,$(words $(GOLIST)),$(GOLIST)))

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

lint:
	[ "$(GOBIN)" = "vgo" ] && $(GOBIN) mod -vendor
	-echo $(GOLIST) | xargs -L1 golint
	@echo
	-deadcode $(GODIR) 2>&1
	@echo
	-find $(GODIR) -type f -exec misspell {} \; $(REVIEWDOG)
	@echo
	-staticcheck $(GOLIST) $(REVIEWDOG)
	@echo
	-errcheck $(GOLIST) $(REVIEWDOG)
	@echo
	-safesql $(GOLIST)
	@echo
	-goconst $(GOLIST) $(REVIEWDOG)
	@echo
	-go vet $(GOLIST) $(REVIEWDOG)
	@echo
	-go vet -shadow $(GOLIST) $(REVIEWDOG)
	@echo
	-aligncheck $(GOLIST) $(REVIEWDOG)
	@echo
	-gosimple $(GOLIST) $(REVIEWDOG)
	@echo
	-unconvert $(GOLIST) $(REVIEWDOG)
	@echo
	-interfacer $(GOLIST) $(REVIEWDOG)

review:
	$(MAKE) lint COMMIT:=master

review-dupl:
	-git diff $(COMMIT) HEAD --name-only --diff-filter=AM|grep -e "\.go$$" | xargs dupl
