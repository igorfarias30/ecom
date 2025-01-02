package cart

import (
	"fmt"

	"github.com/igorfarias30/ecom/types"
)

func getCartItemsID(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))

	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductId)
		}

		productIds[i] = item.ProductId
	}

	return productIds, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// check if all products are actually in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}

	// calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce quantity of products in our db
	for _, item := range items {
		product := productMap[item.ProductId]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}

	// create the order
	for _, item := range items {
		h.store.CreateOrderItems(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductId].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductId]
		if !ok {
			fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductId)
		}

		if product.Quantity < item.Quantity {
			fmt.Errorf("product %d is not available in the quantity requested", item.ProductId)
		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += float64(item.Quantity) * productMap[item.ProductId].Price
	}
	return totalPrice
}
