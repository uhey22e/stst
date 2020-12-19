PGFORMATTER := pg_format
PGFORMATTER_OPTS := -b

.PHONY: all
all: bin/stst lint

bin/stst: ./cmd/stst/main.go
	go build -o $@ ./cmd/stst


.PHONY: lint
lint:
	find testdata -type f -name '*.sql' -exec \
		$(PGFORMATTER) $(PGFORMATTER_OPTS) -o {} {} \;
