package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/igorfarias30/ecom/service/auth"
	"github.com/igorfarias30/ecom/types"
	"github.com/igorfarias30/ecom/utils"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/product", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/product/{id}", h.handleGetProduct).Methods(http.MethodGet)

	// admin routes
	router.HandleFunc("/product", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["id"]

	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	productId, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product id"))
		return
	}

	product, err := h.store.GetProductById(productId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload

	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// check if product already exists
	var productInDb *types.Product
	productInDb, err := h.store.GetProductByName(product.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if productInDb.ID != 0 {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("product with name %v already exists", productInDb.Name))
		return
	}

	err = h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, product)
}
