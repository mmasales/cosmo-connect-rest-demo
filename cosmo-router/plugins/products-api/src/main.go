package main

// This plugin implements the products-api subgraph by making REST calls
// to the demo API at https://demo-api.wundergraph.com/products.
//
// To use your own REST API:
//   1. Change BASE_URL to your API's base URL.
//   2. Add auth headers in newClient() if your API requires them.
//   3. Adjust the JSON field mapping in each method if your schema differs.
//
// The generated pb (protobuf) types in ./generated are produced by running:
//   make generate
//
// Do not edit the generated files — they are regenerated from schema.graphql.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "products-api/src/generated"
)

// BASE_URL is the root of the REST API.
// Override with the PRODUCTS_API_BASE_URL environment variable.
const defaultBaseURL = "https://demo-api.wundergraph.com"

// ---- REST client -------------------------------------------------------
// Centralizes base URL and shared headers — equivalent to Apollo's @source.

type restClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func newClient() *restClient {
	baseURL := os.Getenv("PRODUCTS_API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &restClient{
		baseURL: baseURL,
		apiKey:  os.Getenv("PRODUCTS_API_KEY"), // optional — not required by the demo API
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *restClient) do(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, 0, err
	}
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

// ---- Product type (mirrors the REST API response) ----------------------

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Inventory   int     `json:"inventory"`
}

func (p product) toProto() *pb.Product {
	return &pb.Product{
		Id:          fmt.Sprintf("%d", p.ID),
		Name:        p.Name,
		Price:       p.Price,
		Description: &p.Description,
		Inventory:   int32(p.Inventory),
	}
}

// ---- Service -----------------------------------------------------------

type Service struct {
	pb.UnimplementedProductsApiServer
	client *restClient
}

func NewService() *Service {
	return &Service{client: newClient()}
}

// GET /products

func (s *Service) QueryProducts(ctx context.Context, _ *pb.QueryProductsRequest) (*pb.QueryProductsResponse, error) {
	data, _, err := s.client.do(ctx, http.MethodGet, "/products", nil)
	if err != nil {
		return nil, err
	}

	var products []product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	result := make([]*pb.Product, len(products))
	for i, p := range products {
		result[i] = p.toProto()
	}
	return &pb.QueryProductsResponse{Products: result}, nil
}

// GET /products/{id}

func (s *Service) QueryProduct(ctx context.Context, req *pb.QueryProductRequest) (*pb.QueryProductResponse, error) {
	data, status, err := s.client.do(ctx, http.MethodGet, "/products/"+req.Id, nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return &pb.QueryProductResponse{Product: nil}, nil
	}

	var p product
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &pb.QueryProductResponse{Product: p.toProto()}, nil
}

// POST /products

func (s *Service) MutationCreateProduct(ctx context.Context, req *pb.MutationCreateProductRequest) (*pb.MutationCreateProductResponse, error) {
	body, _ := json.Marshal(map[string]any{
		"name":        req.Input.Name,
		"price":       req.Input.Price,
		"description": req.Input.Description,
		"inventory":   req.Input.Inventory,
	})

	data, _, err := s.client.do(ctx, http.MethodPost, "/products", body)
	if err != nil {
		return nil, err
	}

	var p product
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &pb.MutationCreateProductResponse{Product: p.toProto()}, nil
}

// PUT /products/{id}

func (s *Service) MutationSetProduct(ctx context.Context, req *pb.MutationSetProductRequest) (*pb.MutationSetProductResponse, error) {
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
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &pb.MutationSetProductResponse{Product: p.toProto()}, nil
}

// PATCH /products/{id}

func (s *Service) MutationUpdateProduct(ctx context.Context, req *pb.MutationUpdateProductRequest) (*pb.MutationUpdateProductResponse, error) {
	patch := map[string]any{}
	if req.Input.Name != nil {
		patch["name"] = *req.Input.Name
	}
	if req.Input.Price != nil {
		patch["price"] = *req.Input.Price
	}
	if req.Input.Description != nil {
		patch["description"] = *req.Input.Description
	}
	if req.Input.Inventory != nil {
		patch["inventory"] = *req.Input.Inventory
	}

	body, _ := json.Marshal(patch)
	data, _, err := s.client.do(ctx, http.MethodPatch, "/products/"+req.Id, body)
	if err != nil {
		return nil, err
	}

	var p product
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &pb.MutationUpdateProductResponse{Product: p.toProto()}, nil
}

// DELETE /products/{id}

func (s *Service) MutationDeleteProduct(ctx context.Context, req *pb.MutationDeleteProductRequest) (*pb.MutationDeleteProductResponse, error) {
	_, status, err := s.client.do(ctx, http.MethodDelete, "/products/"+req.Id, nil)
	if err != nil {
		return nil, err
	}
	return &pb.MutationDeleteProductResponse{Success: status == http.StatusNoContent}, nil
}

// ---- Entry point -------------------------------------------------------

func main() {
	svc := NewService()
	log.Println("Starting products-api plugin")
	if err := pb.Serve(svc); err != nil {
		log.Fatal(err)
	}
}
