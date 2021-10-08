PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 windows/amd64/.exe

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

	@echo "==> Building ${os}-${arch} canoegen"
	@$(GO) build -v ${LDFLAGS} -o build/tools/$(os)-$(arch)/canoegen$(ext) github.com/dotstart/canoe/cmd/canoegen
