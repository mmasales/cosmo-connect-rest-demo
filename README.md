# Cosmo Connect — REST API Demo

A minimal, runnable example showing how to integrate a REST API into a Cosmo federated graph using **Cosmo Connect**.

Clone it, run one command, and you have a working federated GraphQL API backed by live REST endpoints — ready to query in the playground.

---

## What this demo does

It exposes a **Products** GraphQL subgraph that calls a live REST API under the hood:

```
https://demo-api.wundergraph.com/products
```

Every GraphQL operation maps to an HTTP method on that API:

| GraphQL operation | HTTP method | Endpoint |
|---|---|---|
| `query { products }` | GET | `/products` |
| `query { product(id: "1") }` | GET | `/products/1` |
| `mutation { createProduct(...) }` | POST | `/products` |
| `mutation { setProduct(...) }` | PUT | `/products/1` |
| `mutation { updateProduct(...) }` | PATCH | `/products/1` |
| `mutation { deleteProduct(...) }` | DELETE | `/products/1` |

---

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [wgc CLI](https://cosmo-docs.wundergraph.com/cli/intro) — `npm install -g wgc@latest`
- [Docker](https://www.docker.com/) (to run the Cosmo Router)

---

## Quickstart

```bash
git clone https://github.com/wundergraph/cosmo-connect-rest-demo.git
cd cosmo-connect-rest-demo
make
```

That's it. `make` will:
1. Generate the gRPC types from `schema.graphql`
2. Build the products-api plugin binary
3. Compose the supergraph into `config.json`
4. Start the Cosmo Router on `http://localhost:3010`

Open **http://localhost:3010** in your browser to access the GraphQL Playground.

---

## Try it

Copy any of these into the playground at `http://localhost:3010`:

```graphql
# List all products
query {
  products {
    id name price inventory
  }
}

# Get one product
query {
  product(id: "1") {
    id name price description
  }
}

# Create a product (POST)
mutation {
  createProduct(input: {
    name: "Wireless Headphones"
    price: 99.99
    inventory: 250
  }) {
    id name
  }
}

# Full update (PUT)
mutation {
  setProduct(id: "1", input: {
    name: "Headphones Pro"
    price: 129.99
    description: "ANC edition"
    inventory: 180
  }) {
    id name price
  }
}

# Partial update (PATCH)
mutation {
  updateProduct(id: "1", input: { price: 109.99 }) {
    id price
  }
}

# Delete (DELETE)
mutation {
  deleteProduct(id: "1")
}
```

More examples are in [`examples/queries.graphql`](./examples/queries.graphql).

---

## Using your own REST API

1. Open `cosmo-router/plugins/products-api/src/schema.graphql` and update the schema to match your API's shape.
2. Open `cosmo-router/plugins/products-api/src/main.go` and change `defaultBaseURL` to your API's base URL.
3. If your API requires auth, add your headers in `newClient()`.
4. Run `make` again.

The `PRODUCTS_API_BASE_URL` and `PRODUCTS_API_KEY` environment variables are also supported if you prefer not to hardcode values.

---

## Project structure

```
cosmo-connect-rest-demo/
├── Makefile                                         # Top-level: make → runs everything
├── examples/
│   └── queries.graphql                              # Ready-to-paste playground queries
└── cosmo-router/
    ├── Makefile                                     # generate → build → compose → start
    ├── config.yaml                                  # Cosmo Router configuration
    ├── graph.yaml                                   # Supergraph composition config
    └── plugins/
        └── products-api/
            ├── Makefile                             # Plugin build steps
            └── src/
                ├── schema.graphql                   # ← Edit this to change your schema
                ├── main.go                          # ← Edit this to change your REST calls
                └── go.mod
```

---

## How it works

Cosmo Connect uses a gRPC adapter layer:

```
GraphQL query
     ↓
Cosmo Router (translates to gRPC)
     ↓
products-api plugin (your Go code)
     ↓
REST API (HTTP)
```

You define the schema, run `make generate` to get type-safe Go stubs, then implement the HTTP calls. The router handles everything else — query planning, batching, and federation.

→ [Full documentation](https://cosmo-docs.wundergraph.com/connect/overview)

---

## License

Apache 2.0 — see [LICENSE](./LICENSE).
