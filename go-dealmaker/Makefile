ABIGEN=go run github.com/ethereum/go-ethereum/cmd/abigen@v1.11.4

contract: 
	$(ABIGEN) --abi abi.json --pkg main --type Contract --out contract.go --bin bytecode.bin
.PHONY: ethereum