.PHONY: clean
clean:
	@rm -f basic/output.txt

.PHONY: test
test: clean
	@flagops
	@cmp basic/output.txt basic/expected.txt --silent || (echo "Basic Test Failed output != expected" && exit 1)
	@echo "All Tests Passed"