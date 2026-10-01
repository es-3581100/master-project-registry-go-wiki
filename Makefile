.PHONY: dev test export clean

dev:
	./scripts/dev.sh

test:
	go test ./...

export:
	./scripts/export.sh

clean:
	rm -rf docs build
