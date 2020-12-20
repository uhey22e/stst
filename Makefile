PGFORMATTER := pg_format
PGFORMATTER_OPTS :=

RMDIR := rm -rf

.PHONY: all
all: bin/stst lint

bin/demo: demo $(wildcard demo/*.go)
	go build -o $@ ./$<

bin/%: cmd/% go.sum $(shell find . -type f -name '*.go')
	go build -o $@ ./$<

.PHONY: test
test:
	go list ./... | grep -v /demo/ | go test -v

.PHONY: lint
lint:
	find . -type f -name '*.sql' -exec \
		$(PGFORMATTER) $(PGFORMATTER_OPTS) -o {} {} \;

.PHONY: clean
clean:
	$(RMDIR) bin
