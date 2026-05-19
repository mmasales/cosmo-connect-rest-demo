# Cosmo Connect — REST API Demo

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

> Add any REST API to your federated graph in minutes using [Cosmo Connect](https://cosmo-docs.wundergraph.com/connect/overview).

[**Quickstart**](#quickstart) · [**How it works**](#how-it-works) · [**All HTTP methods**](#all-http-methods) · [**Use your own API**](#use-your-own-api) · [**Docs**](https://cosmo-docs.wundergraph.com/connect/overview)

---

## Overview

This repo shows how to integrate a REST API into a [Cosmo](https://github.com/wundergraph/cosmo) federated graph using Cosmo Connect. You define a GraphQL schema, implement the HTTP calls in a Go plugin, and the Cosmo Router handles the rest.

All examples use a live demo API:

```
https://cosmo-connect-demo-api-production.up.railway.app/products
```

Try it now:

```bash
curl https://cosmo-connect-demo-api-production.up.railway.app/products
```

---

## Quickstart

**Prerequisites:** [Go 1.22+](https://go.dev/dl/) · [wgc CLI](https://cosmo-docs.wundergraph.com/cli/intro) (`npm i -g wgc@latest`) · [Docker](https://www.docker.com/)

```bash
git clone https://github.com/mmasales/cosmo-connect-rest-demo.git
cd cosmo-connect-rest-demo
make
```

Open **http://localhost:3010** in your browser. That's it.

---

## How it works

```
GraphQL query → Cosmo Router → gRPC → Plugin → REST API
```

1. You define a **GraphQL schema** (`plugins/products-api/src/schema.graphql`)
2. Run `make generate` to get type-safe Go stubs from the schema
3. Implement the **HTTP calls** in `plugins/products-api/src/main.go`
4. Run `make` to build, compose, and start the router

---

## All HTTP methods

Every GraphQL operation maps to an HTTP method:

| GraphQL | HTTP | Endpoint |
|---|---|---|
| `query { products }` | GET | `/products` |
| `query { product(id: "1") }` | GET | `/products/1` |
| `mutation { createProduct(...) }` | POST | `/products` |
| `mutation { setProduct(...) }` | PUT | `/products/1` |
| `mutation { updateProduct(...) }` | PATCH | `/products/1` |
| `mutation { deleteProduct(...) }` | DELETE | `/products/1` |

Copy any of these into the playground at `http://localhost:3010`:

```graphql
# GET /products
query {
  products { id name price inventory }
}

# POST /products
mutation {
  createProduct(input: { name: "Standing Desk", price: 499.99, inventory: 50 }) {
    id name
  }
}

# PATCH /products/1
mutation {
  updateProduct(id: "1", input: { price: 79.99 }) {
    id price
  }
}

# DELETE /products/1
mutation {
  deleteProduct(id: "1")
}
```

More examples in [`examples/queries.graphql`](./examples/queries.graphql).

---

## Use your own API

1. Update the schema in `plugins/products-api/src/schema.graphql`
2. Run `make generate` to regenerate the Go stubs
3. Update the REST calls in `plugins/products-api/src/main.go`
4. Set your base URL: `PRODUCTS_API_BASE_URL=https://your-api.com make`

---

## Project structure

```
cosmo-connect-rest-demo/
├── Makefile                                    # build + compose + start
├── config.yaml                                 # router config
├── graph.yaml                                  # supergraph composition
├── examples/queries.graphql                    # ready-to-paste queries
└── plugins/products-api/
    ├── src/
    │   ├── schema.graphql                      # ← your GraphQL schema
    │   └── main.go                             # ← your REST calls
    └── generated/                              # auto-generated, do not edit
```

---

## Further reading

- [Cosmo Connect Overview](https://cosmo-docs.wundergraph.com/connect/overview)
- [Router Plugins](https://cosmo-docs.wundergraph.com/router/gRPC/plugins)
- [Full Cosmo Docs](https://cosmo-docs.wundergraph.com)

## License

Apache 2.0 — see [LICENSE](./LICENSE).
