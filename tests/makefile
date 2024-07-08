
SUBDIRS := $(wildcard */)

export string_feature = {"defaultVariant":"string","variants":[{"name":"string","targetingKey":"","criteria":[],"value":"test string"}]}
export integer_feature = {"defaultVariant":"int","variants":[{"name":"int","targetingKey":"","criteria":[],"value":3}]}
export bool_feature = {"defaultVariant":"bool","variants":[{"name":"bool","targetingKey":"","criteria":[],"value":true}]}

.PHONY: test
test: ${SUBDIRS}
	@./flagops

	@for f in $^; do \
		echo "Running [$${f}] test" && \
		$(MAKE) -C $${f}; \
	done

	@echo "All Tests Passed"

.PHONY: clean
clean:
	@rm -f basic/output.txt