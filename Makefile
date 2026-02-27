.PHONY: build-node build-client build clean

build-node:
	@echo "Building Go Node..."
	cd node && go build -o ../node-bin

build-client:
	@echo "Building Go Client CLI..."
	cd client && go build -o ../client-bin

build: build-node build-client
	@echo "All binaries built successfully."

clean:
	@echo "Cleaning up compiled binaries..."
	rm -f node-bin client-bin
	rm -f *.enc *.json *.txt
	rm -rf target/ contract/target/
	@echo "Clean complete."
