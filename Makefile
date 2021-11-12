APPLICATION_VERSION := 0.1.0
APPLICATION_COMMIT_HASH := `git log -1 --pretty=format:"%H"`
APPLICATION_TIMESTAMP := `date --utc "+%s"`

LDFLAGS :=-X 'github.com/dotstart/githoard/internal.version=${APPLICATION_VERSION}' -X 'github.com/dotstart/githoard/internal.commitHash=${APPLICATION_COMMIT_HASH}' -X 'github.com/dotstart/githoard/internal.timestampRaw=${APPLICATION_TIMESTAMP}'

GO := $(shell command -v go 2> /dev/null)
TAR := $(shell command -v tar 2> /dev/null)
export

PLATFORMS := darwin/amd64 linux/amd64 linux/arm windows/amd64/.exe

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

all: check-env $(PLATFORMS)

check-env:
	@echo "==> Checking prerequisites"
	@echo -n "Checking for go ... "
ifndef GO
	@echo "Not Found"
	$(error "go is unavailable")
endif
	@echo $(GO)
	@echo -n "Checking for tar ... "
ifndef TAR
	@echo "Not Found"
	$(error "tar is unavailable")
endif
	@echo $(GO)
	@echo ""

clean:
	@echo "==> Clearing previous build data"
	@rm -rf build/ || true
	@$(GO) clean -cache

.ONESHELL:
$(PLATFORMS):
	@export GOOS=$(os);
	@export GOARCH=$(arch);

	@echo "==> Building ${os}-${arch}"
	@$(GO) build -v -ldflags "${LDFLAGS}" -o build/$(os)-$(arch)/githoard$(ext) github.com/dotstart/githoard/cmd/githoard
	@$(TAR) -C "build/$(os)-$(arch)/" -czvf "build/githoard_$(os)-$(arch).tar.gz" "githoard$(ext)"

.PHONY: all
