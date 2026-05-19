.PHONY: start build compose generate clean

# Clone, run `make` — that's it.
start: build compose
	$(MAKE) -C cosmo-router start

build:
	$(MAKE) -C cosmo-router build

compose:
	$(MAKE) -C cosmo-router compose

generate:
	$(MAKE) -C cosmo-router generate

clean:
	$(MAKE) -C cosmo-router clean
