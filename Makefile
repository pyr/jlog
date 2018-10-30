PKG = 		jlog
MAIN = 		$(PKG).go
RM =		rm -f

.PHONY: all
all: $(PKG)

$(PKG): $(PKG).go
	go build -o $(PKG)

.PHONY: clean
clean:
	$(RM) $(PKG)

.PHONY: version
version:
	@echo $(VERSION)
