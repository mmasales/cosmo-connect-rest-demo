.PHONY: start build compose generate clean

ROUTER_VERSION ?= latest
PLUGIN_NAME    := products-api
PLUGIN_DIR     := plugins/$(PLUGIN_NAME)

start: build compose
	@echo "→ Starting Cosmo Router on http://localhost:3010"
	docker run \
		--name cosmo-router \
		--rm \
		-p 3010:3002 \
		--add-host=host.docker.internal:host-gateway \
		--platform=linux/amd64 \
		-e DEV_MODE=true \
		-e LISTEN_ADDR=0.0.0.0:3002 \
		-e EXECUTION_CONFIG_FILE_PATH=/config/config.json \
		-v "$(PWD)/config.yaml:/config.yaml" \
		-v "$(PWD)/config.json:/config/config.json" \
		-v "$(PWD)/plugins:/plugins" \
		ghcr.io/wundergraph/cosmo/router:$(ROUTER_VERSION)

build: generate
	@echo "→ Building plugin: $(PLUGIN_NAME)"
	$(MAKE) -C $(PLUGIN_DIR) build

generate:
	@echo "→ Generating proto from schema.graphql"
	$(MAKE) -C $(PLUGIN_DIR) generate

compose:
	@echo "→ Composing supergraph → config.json"
	wgc router compose --input graph.yaml --out config.json

clean:
	$(MAKE) -C $(PLUGIN_DIR) clean
	rm -f config.json

	