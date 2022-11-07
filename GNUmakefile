default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: TF_ACC=1
testacc:
	gotestsum --no-color=false -ftestname -- -race ./...
