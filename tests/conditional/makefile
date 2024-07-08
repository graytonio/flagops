.PHONY: check
check:
	@(cmp output.txt expected.txt --silent && echo "Conditional Test Passed") || (echo "Conditional Test Failed output != expected" && exit 1)