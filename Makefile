PGFORMATTER := pg_format
PGFORMATTER_OPTS := -b


.PHONY: lint
lint:
	find testdata -type f -name '*.sql' -exec \
		$(PGFORMATTER) $(PGFORMATTER_OPTS) -o {} {} \;
