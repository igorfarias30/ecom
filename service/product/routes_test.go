package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/igorfarias30/ecom/types"
)

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(productStore, userStore)

	t.Run("should handle get products", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail creating a product if the payload is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail creating product when the product with name already exists", func(t *testing.T) {
		alreadyExistedProduct := "super guitar"
		newProductPayload := types.CreateProductPayload{
			Name:        alreadyExistedProduct,
			Price:       100,
			Image:       "test.jpg",
			Description: "test description",
			Quantity:    5,
		}

		marshalled, err := json.Marshal(newProductPayload)
		if err != nil {
			t.Fatal(err)
		}

		productStore.MockGetProductByName = func(name string) (*types.Product, error) {
			if name == alreadyExistedProduct {
				return &types.Product{ID: 1, Name: alreadyExistedProduct}, nil
			}
			return nil, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/product", handler.handleCreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("expected status code %d, got %d", http.StatusConflict, rr.Code)
		}
	})

	t.Run("should handle creating a product", func(t *testing.T) {
		newProductPayload := types.CreateProductPayload{
			Name:        "Guitar Test",
			Price:       100,
			Image:       "test.jpg",
			Description: "test description",
			Quantity:    5,
		}

		marshalled, err := json.Marshal(newProductPayload)
		if err != nil {
			t.Fatal(err)
		}

		productStore.MockGetProductByName = func(name string) (*types.Product, error) {
			return &types.Product{}, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockProductStore struct {
	MockGetProductByName func(name string) (*types.Product, error)
}

func (m *mockProductStore) GetProducts() ([]*types.Product, error) {
	return []*types.Product{}, nil
}

func (m *mockProductStore) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	return nil, nil
}

func (m *mockProductStore) GetProductByName(name string) (*types.Product, error) {
	return m.MockGetProductByName(name)
}

func (m *mockProductStore) CreateProduct(product types.CreateProductPayload) error {
	return nil
}

func (m *mockProductStore) UpdateProduct(product types.Product) error {
	return nil
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
