.PHONY: conformance docs safety

conformance:
	./scripts/run-test262.sh all

docs:
	./scripts/check-docs.sh

safety:
	./scripts/check-safety.sh
