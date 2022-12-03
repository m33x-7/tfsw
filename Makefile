.DEFAULT_GOAL := tfsw
CMD = GOARCH=$(GOARCH) GOOS=$(GOOS) go build $(GO_BUILD_ARGS) -o $@
GO_BUILD_ARGS = -v -ldflags "-s -w -X tfsw/cmd.buildVer=$(VERSION)"
GOARCH = $(word 3, $(TARGET_WORDS))
GOOS = $(word 2, $(TARGET_WORDS))
TARGET_WORDS = $(subst _, ,$(subst .,_,$@))
VERSION ?= $(shell git rev-parse --short HEAD)


.PHONY: all clean install

all: tfsw_darwin_amd64.gz \
	tfsw_darwin_arm64.gz \
	tfsw_freebsd_amd64.gz \
	tfsw_linux_amd64.gz \
	tfsw_windows_amd64.zip

clean:
	-rm -vf tfsw
	-rm -vf tfsw.exe
	-rm -vf tfsw_darwin_*
	-rm -vf tfsw_freebsd_*
	-rm -vf tfsw_linux_*
	-rm -vf tfsw_windows_*
	-rm -vf tfsw_SHA256SUMS

install:
	@echo Installing tfsw ...
	mkdir -p ~/bin
	cp tfsw ~/bin/tfsw

tfsw: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	go build $(GO_BUILD_ARGS)

tfsw_darwin_amd64: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	$(CMD)

tfsw_darwin_amd64.gz: tfsw_darwin_amd64
	cp tfsw_darwin_amd64 tfsw
	gzip -Nc tfsw >> tfsw_darwin_amd64.gz
	rm tfsw

tfsw_darwin_arm64: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	$(CMD)

tfsw_darwin_arm64.gz: tfsw_darwin_arm64
	cp tfsw_darwin_amd64 tfsw
	gzip -Nc tfsw >> tfsw_darwin_arm64.gz
	rm tfsw

tfsw_freebsd_amd64: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	$(CMD)

tfsw_freebsd_amd64.gz: tfsw_freebsd_amd64
	cp tfsw_freebsd_amd64 tfsw
	gzip -Nc tfsw >> tfsw_freebsd_amd64.gz
	rm tfsw

tfsw_linux_amd64: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	$(CMD)

tfsw_linux_amd64.gz: tfsw_linux_amd64
	cp tfsw_linux_amd64 tfsw
	gzip -Nc tfsw >> tfsw_linux_amd64.gz
	rm tfsw

tfsw_SHA256SUMS: tfsw_darwin_amd64.gz \
	tfsw_darwin_arm64.gz \
	tfsw_freebsd_amd64.gz \
	tfsw_linux_amd64.gz \
	tfsw_windows_amd64.zip
	@echo Creating SHA256SUMS ...
	sha256sum tfsw_darwin_amd64.gz \
		tfsw_darwin_arm64.gz \
		tfsw_freebsd_amd64.gz \
		tfsw_linux_amd64.gz \
		tfsw_windows_amd64.zip \
		> tfsw_SHA256SUMS
	sha256sum -c tfsw_SHA256SUMS

tfsw_SHA256SUMS.gpg: tfsw_SHA256SUMS
	rm tfsw_SHA256SUMS.gpg
	gpg -s tfsw_SHA256SUMS

tfsw_windows_amd64.exe: cmd/*.go internal/**/*.go
	@echo Building $@ ...
	$(CMD)

tfsw_windows_amd64.zip: tfsw_windows_amd64.exe
	cp tfsw_windows_amd64.exe tfsw.exe
	zip tfsw_windows_amd64.zip tfsw.exe
	rm tfsw.exe
