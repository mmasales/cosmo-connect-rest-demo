# ProductsApi Plugin - Cosmo Router Example

This repository contains a simple Cosmo Router plugin that showcases how to design APIs with GraphQL Federation but implement them using gRPC methods instead of traditional resolvers.

## What is this demo about?

This demo illustrates a key pattern in Cosmo Router plugin development:
- **Design with GraphQL**: Define your API using GraphQL schema
- **Implement with gRPC**: Instead of writing GraphQL resolvers, implement gRPC service methods
- **Bridge the gap**: The Cosmo router connects GraphQL operations to your gRPC implementations
- **Test-Driven Development**: Test your gRPC service implementation with gRPC client and server without external dependencies

The plugin demonstrates:
- How GraphQL types and operations map to gRPC RPC methods
- Simple "Hello World" implementation
- Proper structure for a Cosmo Router plugin
- How to test your gRPC implementation with gRPC client and server without external dependencies

## Getting Started

Plugin structure:

   ```
    plugins/products-api/
    ├── go.mod                # Go module file with dependencies
    ├── go.sum                # Go checksums file
    ├── src/
    │   ├── main.go           # Main plugin implementation
    │   ├── main_test.go      # Tests for the plugin
    │   └── schema.graphql    # GraphQL schema defining the API
    ├── generated/            # Generated code (created during build)
    └── bin/                  # Compiled binaries (created during build)
        └── plugin            # The compiled plugin binary
   ```

## 🔧 Customizing Your Plugin

- Change the GraphQL schema in `src/schema.graphql` and regenerate the code with `make generate`.
- Implement the changes in `src/main.go` and test your implementation with `make test`.
- Compose your supergraph with `make compose` and restart the router with `make start`.

## 📚 Learn More

For more information about Cosmo and building router plugins:
- [Cosmo Documentation](https://cosmo-docs.wundergraph.com/)
- [Cosmo Router Plugins Guide](https://cosmo-docs.wundergraph.com/connect/plugins)

---

<p align="center">Made with ❤️ by <a href="https://wundergraph.com">WunderGraph</a></p>