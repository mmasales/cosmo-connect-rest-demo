.PHONY: generate build clean

PLUGIN_NAME := products-api
SRC_DIR     := src
OUT_DIR     := bin

generate:
	@echo "  → Generating proto from schema.graphql"
	wgc router plugin generate --schema $(SRC_DIR)/schema.graphql --out $(SRC_DIR)/generated

build: generate
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(PLUGIN_NAME) ./$(SRC_DIR)/...
	cp $(OUT_DIR)/$(PLUGIN_NAME) ./$(PLUGIN_NAME)

clean:
	rm -rf $(OUT_DIR) $(SRC_DIR)/generated $(PLUGIN_NAME)