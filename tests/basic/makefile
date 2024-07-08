.PHONY: check
check:
	@(cmp output.txt expected.txt --silent && echo "Basic Test Passed") || (echo "Basic Test Failed output != expected" && exit 1)