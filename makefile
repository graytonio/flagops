build:
	goreleaser build --snapshot --clean --single-target

build-test:
	goreleaser build --snapshot --clean --single-target --output tests/flagops
	$(MAKE) -C tests