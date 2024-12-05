
.PHONY: debug
debug:
	dlv debug --headless --listen=:2345 . -- $(CMD)

build:
	mkdir target
	go build -o target/vaultview -ldflags=" -X 'vaultview/pkg/models.version=dev' "
