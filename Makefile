VERSION = 	v0.1.1-snapshot
PKG = 		jlog
MAIN = 		$(PKG).go
RM =		rm -f
DEP =		$(GOPATH)/bin/dep

.PHONY: all
all: $(PKG)

$(PKG): $(PKG).go
	go build

.PHONY: clean
clean:
	$(RM) $(PKG)

.PHONY: version
version:
	@echo $(VERSION)
