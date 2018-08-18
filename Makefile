PROJECT_SINCE:=$(shell git log --pretty=format:"%ad" --date=unix|tail -1)
AUTO_COUNT_SINCE:=$(shell echo $$(((`date -u +%s`-$(PROJECT_SINCE))/(24*60*60))))
AUTO_COUNT_LOG:=$(shell git log --since=midnight --oneline|wc -l|tr -d " ")
COMMIT:=4b825dc
REVIEWDOG:=| reviewdog -efm='%f:%l:%c: %m' -diff="git diff $(COMMIT) HEAD"

GOBIN=go1.11rc1
VGOBIN:=GO111MODULE=on go1.11rc1
PKG:=$(shell $(GOBIN) list)
NAME:=$(notdir $(PKG))
GOLIST:=$(shell $(GOBIN) list ./...)
GODIR:=$(patsubst $(PKG)/%,%,$(wordlist 2,$(words $(GOLIST)),$(GOLIST)))

GO:=$(find . -type d -name vendor -prune -or -type f -name "*.go" -print)
LIBGO:=$(wildcard lib/*.go)
LIB:=$(LIBGO:.go=.so)
VENDOR:=vendor/golang.org

.SUFFIXES: .go .so
.go.so:
	$(GOBIN) build -buildmode=c-shared -o $@ $<

all: $(VENDOR) $(GO) $(LIB) test
	$(GOBIN) build -ldflags=" \
-X main.serial=$(AUTO_COUNT_SINCE).$(AUTO_COUNT_LOG) \
-X main.hash=$(shell git describe --always --dirty=+) \
-X \"main.build=$(shell LANG=en date -u +'%b %d %T %Y')\" \
"

go.mod:
	touch $@

$(VENDOR): go.mod
	$(VGOBIN) mod vendor
	git checkout vendor/v

clean:
	rm -f $(NAME) $(wildcard lib/*.h) $(wildcard lib/*.so) vendor/modules.txt
	find vendor -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;

test:
	$(GOBIN) test

lint: $(VENDOR)
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
	-$(GOBIN) vet $(GOLIST) $(REVIEWDOG)
	@echo
	-$(GOBIN) vet -shadow $(GOLIST) $(REVIEWDOG)
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
