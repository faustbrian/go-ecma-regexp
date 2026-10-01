.PHONY: conformance docs safety interoperability

conformance:
	./scripts/run-test262.sh all

docs:
	./scripts/check-docs.sh

safety:
	./scripts/check-safety.sh

interoperability:
	./scripts/run-differential.sh
