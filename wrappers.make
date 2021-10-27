PLATFORMS := darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386/.exe windows/amd64/.exe

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

build: $(PLATFORMS)

.ONESHELL:
$(PLATFORMS):
	@export GOOS=$(os);
	@export GOARCH=$(arch);

	@echo "==> Building ${os}-${arch} CLI wrapper"
	@$(GO) build -v -ldflags "${LDFLAGS}" -o build/wrappers/$(os)-$(arch)/canoew$(ext) github.com/dotstart/canoe/cmd/canoew-cli

ifdef UPX_BIN
	@echo "==> Compressing ${os}-${arch} CLI wrapper with UPX"
	@$(UPX_BIN) ${UPX_FLAGS} build/wrappers/$(os)-$(arch)/canoew$(ext)
endif

	@if [ "$(os)" = "windows" ]; then\
		echo "==> Building ${os}-${arch} GUI wrapper"; \
		$(GO) build -v -ldflags "-H=windowsgui ${LDFLAGS}" -o build/wrappers/$(os)-$(arch)/canoew-gui$(ext) github.com/dotstart/canoe/cmd/canoew-gui; \

ifdef UPX_BIN
		@echo "==> Compressing ${os}-${arch} GUI wrapper with UPX"
		@$(UPX_BIN) ${UPX_FLAGS} build/wrappers/$(os)-$(arch)/canoew-gui$(ext)
endif
	fi
