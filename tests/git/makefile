.PHONY: check
check:
	@curl -o output.txt https://raw.githubusercontent.com/graytonio/flagops-test-repo/main/git-test/test.txt
	@(cmp output.txt expected.txt --silent && echo "Env Test Passed") || (echo "Env Test Failed output != expected" && exit 1)
	@curl -o output.txt https://raw.githubusercontent.com/graytonio/flagops-test-repo/test-branch/git-test/test.txt
	@(cmp output.txt expected.txt --silent && echo "Env Test Passed") || (echo "Env Test Failed output != expected" && exit 1)