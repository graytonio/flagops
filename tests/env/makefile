.PHONY: check
check:
	@(cmp output.txt expected.txt --silent && echo "Env Test Passed") || (echo "Env Test Failed output != expected" && exit 1)