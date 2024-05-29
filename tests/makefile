
SUBDIRS := $(wildcard */)

.PHONY: test
test: ${SUBDIRS}
	@./flagops

	@for f in $^; do \
		echo "Running [$${f}] test" && \
		cmp $${f}/output.txt $${f}/expected.txt --silent || (echo "$${f} Test Failed output != expected" && exit 1); \
	done

	@echo "All Tests Passed"

.PHONY: clean
clean:
	@rm -f basic/output.txt