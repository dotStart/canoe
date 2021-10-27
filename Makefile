APPLICATION_VERSION := 0.1.0

LDFLAGS := -w -s -X github.com/dotstart/canoe/internal.version=${APPLICATION_VERSION}
UPX_FLAGS := --ultra-brute

GO := $(shell command -v go 2> /dev/null)
UPX_BIN := $(shell command -v upx 2> /dev/null)
export

all: check-env tools package

clean:
	@echo "==> Clearing previous build data"
	@rm -rf build/wrappers/ || true
	@rm -rf build/tools/ || true
	@$(GO) clean -cache

check-env:
	@echo "==> Checking prerequisites"
	@echo -n "Checking for go ... "
ifndef GO
	@echo "Not Found"
	$(error "go is unavailable")
endif
	@echo $(GO)
	@echo -n "Checking for upx ... "
ifndef UPX_BIN
		@echo "Not found - Builds may be larger than usual"
else
		@echo $(UPX_BIN)
endif

tools: wrappers
	$(MAKE) -f tools.make build

wrappers:
	$(MAKE) -f wrappers.make build

package:
	@echo "==> Creating distribution packages"
	@for dir in build/tools/*; do if [ -d "$$dir" ]; then tar -czvf "$(basename "$$dir").tar.gz" --xform="s,$$dir/,," "$$dir"; fi; done
	@echo ""

.PHONY: all
