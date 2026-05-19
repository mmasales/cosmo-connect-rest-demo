# REST API Integration

Connect any REST API to your federated graph using Cosmo Connect. Define a GraphQL schema, implement the HTTP calls in a plugin, and the Cosmo Router handles query planning, batching, and federation.

> **Try it now:** Clone the [working example](https://github.com/mmasales/cosmo-connect-rest-demo) and run `make` to see all HTTP methods working against a live API.

---

## Setup

```bash
# Scaffold a new project
wgc router plugin init products-api -p myproject --language go
cd myproject

# Edit your schema, then generate Go stubs
cd plugins/products-api && make generate

# Build, compose, and start the router
cd ../.. && make
```

Open `http://localhost:3010` to query your graph.

---

## GET

Expose a GET endpoint with a `Query` field.

```graphql
type Query {
  products: [Product!]!
  product(id: ID!): Product
}

type Product @key(fields: "id") {
  id: ID!
  name: String!
  price: Float!
  description: String
  inventory: Int!
}
```

```go
func (s *ProductsApiService) QueryProducts(ctx context.Context, _ *service.QueryProductsRequest) (*service.QueryProductsResponse, error) {
    data, _, err := s.client.do(ctx, http.MethodGet, "/products", nil)
    if err != nil {
        return nil, err
    }
    var products []product
    json.Unmarshal(data, &products)
    result := make([]*service.Product, len(products))
    for i, p := range products {
        result[i] = p.toProto()
    }
    return &service.QueryProductsResponse{Products: result}, nil
}

func (s *ProductsApiService) QueryProduct(ctx context.Context, req *service.QueryProductRequest) (*service.QueryProductResponse, error) {
    data, status, err := s.client.do(ctx, http.MethodGet, "/products/"+req.Id, nil)
    if err != nil {
        return nil, err
    }
    if status == http.StatusNotFound {
        return &service.QueryProductResponse{}, nil
    }
    var p product
    json.Unmarshal(data, &p)
    return &service.QueryProductResponse{Product: p.toProto()}, nil
}
```

```graphql
query {
  products { id name price inventory }
}

query {
  product(id: "1") { id name price description }
}
```

---

## POST

Expose a POST endpoint with a `Mutation` field.

```graphql
type Mutation {
  createProduct(input: CreateProductInput!): Product!
}

input CreateProductInput {
  name: String!
  price: Float!
  description: String
  inventory: Int!
}
```

```go
func (s *ProductsApiService) MutationCreateProduct(ctx context.Context, req *service.MutationCreateProductRequest) (*service.MutationCreateProductResponse, error) {
    body, _ := json.Marshal(map[string]any{
        "name":        req.Input.Name,
        "price":       req.Input.Price,
        "description": req.Input.Description.GetValue(),
        "inventory":   req.Input.Inventory,
    })
    data, _, err := s.client.do(ctx, http.MethodPost, "/products", body)
    if err != nil {
        return nil, err
    }
    var p product
    json.Unmarshal(data, &p)
    return &service.MutationCreateProductResponse{CreateProduct: p.toProto()}, nil
}
```

```graphql
mutation {
  createProduct(input: {
    name: "Standing Desk"
    price: 499.99
    description: "Electric height adjustable"
    inventory: 50
  }) {
    id name price
  }
}
```

---

## PUT

Use PUT for a full replacement of a resource.

```graphql
type Mutation {
  setProduct(id: ID!, input: SetProductInput!): Product!
}

input SetProductInput {
  name: String!
  price: Float!
  description: String!
  inventory: Int!
}
```

```go
func (s *ProductsApiService) MutationSetProduct(ctx context.Context, req *service.MutationSetProductRequest) (*service.MutationSetProductResponse, error) {
    body, _ := json.Marshal(map[string]any{
        "name":        req.Input.Name,
        "price":       req.Input.Price,
        "description": req.Input.Description,
        "inventory":   req.Input.Inventory,
    })
    data, _, err := s.client.do(ctx, http.MethodPut, "/products/"+req.Id, body)
    if err != nil {
        return nil, err
    }
    var p product
    json.Unmarshal(data, &p)
    return &service.MutationSetProductResponse{SetProduct: p.toProto()}, nil
}
```

```graphql
mutation {
  setProduct(id: "1", input: {
    name: "Wireless Headphones Pro"
    price: 129.99
    description: "Over-ear ANC edition"
    inventory: 180
  }) {
    id name price
  }
}
```

---

## PATCH

Use PATCH to update only specific fields. Optional fields in GraphQL map to `google.protobuf.StringValue` / `DoubleValue` / `Int32Value` in the generated proto — use `.GetValue()` to safely unwrap them.

```graphql
type Mutation {
  updateProduct(id: ID!, input: UpdateProductInput!): Product!
}

input UpdateProductInput {
  name: String
  price: Float
  description: String
  inventory: Int
}
```

```go
func (s *ProductsApiService) MutationUpdateProduct(ctx context.Context, req *service.MutationUpdateProductRequest) (*service.MutationUpdateProductResponse, error) {
    patch := map[string]any{}
    if req.Input.Name != nil        { patch["name"] = req.Input.Name.GetValue() }
    if req.Input.Price != nil       { patch["price"] = req.Input.Price.GetValue() }
    if req.Input.Description != nil { patch["description"] = req.Input.Description.GetValue() }
    if req.Input.Inventory != nil   { patch["inventory"] = req.Input.Inventory.GetValue() }

    body, _ := json.Marshal(patch)
    data, _, err := s.client.do(ctx, http.MethodPatch, "/products/"+req.Id, body)
    if err != nil {
        return nil, err
    }
    var p product
    json.Unmarshal(data, &p)
    return &service.MutationUpdateProductResponse{UpdateProduct: p.toProto()}, nil
}
```

```graphql
mutation {
  updateProduct(id: "1", input: { price: 79.99 }) {
    id price
  }
}
```

---

## DELETE

```graphql
type Mutation {
  deleteProduct(id: ID!): Boolean!
}
```

```go
func (s *ProductsApiService) MutationDeleteProduct(ctx context.Context, req *service.MutationDeleteProductRequest) (*service.MutationDeleteProductResponse, error) {
    _, status, err := s.client.do(ctx, http.MethodDelete, "/products/"+req.Id, nil)
    if err != nil {
        return nil, err
    }
    return &service.MutationDeleteProductResponse{DeleteProduct: status == http.StatusNoContent}, nil
}
```

```graphql
mutation {
  deleteProduct(id: "1")
}
```

---

## Sharing base URL and headers

Centralize your base URL and auth headers in a shared client — the equivalent of Apollo's `@source` directive.

```go
type restClient struct {
    baseURL string
    apiKey  string
    http    *http.Client
}

func newClient() *restClient {
    baseURL := os.Getenv("PRODUCTS_API_BASE_URL")
    if baseURL == "" {
        baseURL = "https://your-api.com"
    }
    return &restClient{
        baseURL: baseURL,
        apiKey:  os.Getenv("PRODUCTS_API_KEY"),
        http:    &http.Client{Timeout: 10 * time.Second},
    }
}

func (c *restClient) do(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
    req, _ := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    if c.apiKey != "" {
        req.Header.Set("X-API-Key", c.apiKey)
    }
    resp, err := c.http.Do(req)
    if err != nil {
        return nil, 0, err
    }
    defer resp.Body.Close()
    data, err := io.ReadAll(resp.Body)
    return data, resp.StatusCode, err
}
```

Pass `newClient()` into your service at startup and reuse it across all resolver methods.

---

## Next steps

- [Working example repo](https://github.com/mmasales/cosmo-connect-rest-demo) 
- [Router Plugins reference](https://cosmo-docs.wundergraph.com/router/gRPC/plugins)
- [gRPC concepts](https://cosmo-docs.wundergraph.com/router/gRPC/concepts)
- [Publishing plugins to Cosmo Cloud](https://cosmo-docs.wundergraph.com/cli/router/plugin/publish)
