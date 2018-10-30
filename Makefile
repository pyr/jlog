PKG = 		jlog
MAIN = 		$(PKG).go
RM =		rm -f

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
