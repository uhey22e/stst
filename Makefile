PGFORMATTER := pg_format
PGFORMATTER_OPTS := -b


.PHONY: lint
lint: testdata/tables/V0_1_0__create_tables.sql
	$(PGFORMATTER) $(PGFORMATTER_OPTS) -o $< $<
