PGFORMATTER := pg_format
PGFORMATTER_OPTS :=

RMDIR := rm -rf

.PHONY: all
all: bin/stst lint

bin/demo: demo $(wildcard demo/*.go)
	go build -o $@ ./$<

bin/%: cmd/% $(wildcard cmd/%/*.go)
	go build -o $@ ./$<

.PHONY: lint
lint: $(wildcard testdata/*.sql)
	find testdata -type f -name '*.sql' -exec \
		$(PGFORMATTER) $(PGFORMATTER_OPTS) -o {} {} \;

.PHONY: clean
clean:
	$(RMDIR) bin
