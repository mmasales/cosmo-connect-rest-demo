package main

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

	service "github.com/wundergraph/cosmo/plugin/generated"
	routerplugin "github.com/wundergraph/cosmo/router-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const defaultBaseURL = "https://cosmo-connect-demo-api-production.up.railway.app"

// ---- REST client -------------------------------------------------------

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
		apiKey:  os.Getenv("PRODUCTS_API_KEY"),
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

// ---- Product type ------------------------------------------------------

type product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Inventory   int     `json:"inventory"`
}

func (p product) toProto() *service.Product {
	return &service.Product{
		Id:          fmt.Sprintf("%d", p.ID),
		Name:        p.Name,
		Price:       p.Price,
		Description: wrapperspb.String(p.Description),
		Inventory:   int32(p.Inventory),
	}
}

// ---- Service -----------------------------------------------------------

type ProductsApiService struct {
	service.UnimplementedProductsApiServiceServer
	client *restClient
}

// GET /products

// GET /products
func (s *ProductsApiService) QueryProducts(ctx context.Context, _ *service.QueryProductsRequest) (*service.QueryProductsResponse, error) {
	data, _, err := s.client.do(ctx, http.MethodGet, "/products", nil)
	if err != nil {
		return nil, err
	}

	// DummyJSON wraps products in {"products": [...]}
	// Fall back to plain array if that fails
	var wrapped struct {
		Products []product `json:"products"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && wrapped.Products != nil {
		result := make([]*service.Product, len(wrapped.Products))
		for i, p := range wrapped.Products {
			result[i] = p.toProto()
		}
		return &service.QueryProductsResponse{Products: result}, nil
	}

	// Plain array (our demo API)
	var products []product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	result := make([]*service.Product, len(products))
	for i, p := range products {
		result[i] = p.toProto()
	}
	return &service.QueryProductsResponse{Products: result}, nil
}

// GET /products/{id}
func (s *ProductsApiService) QueryProduct(ctx context.Context, req *service.QueryProductRequest) (*service.QueryProductResponse, error) {
	data, status, err := s.client.do(ctx, http.MethodGet, "/products/"+req.Id, nil)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return &service.QueryProductResponse{}, nil
	}
	var p product
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &service.QueryProductResponse{Product: p.toProto()}, nil
}

// LookupProductById — entity resolver for federation
func (s *ProductsApiService) LookupProductById(ctx context.Context, req *service.LookupProductByIdRequest) (*service.LookupProductByIdResponse, error) {
	result := make([]*service.Product, len(req.Keys))
	for i, key := range req.Keys {
		data, status, err := s.client.do(ctx, http.MethodGet, "/products/"+key.Id, nil)
		if err != nil {
			return nil, err
		}
		if status == http.StatusNotFound {
			result[i] = nil
			continue
		}
		var p product
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, err
		}
		result[i] = p.toProto()
	}
	return &service.LookupProductByIdResponse{Result: result}, nil
}

// POST /products
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
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &service.MutationCreateProductResponse{CreateProduct: p.toProto()}, nil
}

// PUT /products/{id}
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
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &service.MutationSetProductResponse{SetProduct: p.toProto()}, nil
}

// PATCH /products/{id}
func (s *ProductsApiService) MutationUpdateProduct(ctx context.Context, req *service.MutationUpdateProductRequest) (*service.MutationUpdateProductResponse, error) {
	patch := map[string]any{}
	if req.Input.Name != nil {
		patch["name"] = req.Input.Name.GetValue()
	}
	if req.Input.Price != nil {
		patch["price"] = req.Input.Price.GetValue()
	}
	if req.Input.Description != nil {
		patch["description"] = req.Input.Description.GetValue()
	}
	if req.Input.Inventory != nil {
		patch["inventory"] = req.Input.Inventory.GetValue()
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
	return &service.MutationUpdateProductResponse{UpdateProduct: p.toProto()}, nil
}

// DELETE /products/{id}
func (s *ProductsApiService) MutationDeleteProduct(ctx context.Context, req *service.MutationDeleteProductRequest) (*service.MutationDeleteProductResponse, error) {
	_, status, err := s.client.do(ctx, http.MethodDelete, "/products/"+req.Id, nil)
	if err != nil {
		return nil, err
	}
	return &service.MutationDeleteProductResponse{DeleteProduct: status == http.StatusNoContent}, nil
}

// ---- Entry point -------------------------------------------------------

func main() {
	client := newClient()
	pl, err := routerplugin.NewRouterPlugin(func(s *grpc.Server) {
		s.RegisterService(&service.ProductsApiService_ServiceDesc, &ProductsApiService{
			client: client,
		})
	}, routerplugin.WithTracing())
	if err != nil {
		log.Fatalf("failed to create router plugin: %v", err)
	}
	pl.Serve()
}
