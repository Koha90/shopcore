package domain

import (
	"errors"
)

// CartStatus represents cart lifecycle state.
type CartStatus string

const (
	CartStatusActive     CartStatus = "active"
	CartStatusCheckedOut CartStatus = "checked_out"
)

var (
	ErrInvalidUserID      error = errors.New("invalid user id")
	ErrCartNotActive      error = errors.New("cart is not active")
	ErrItemNotFound       error = errors.New("item not found")
	ErrItemAlreadyExists  error = errors.New("item already exists in cart")
	ErrInvalidItemQuality error = errors.New("invalid item quantity")
	ErrCartEmpty          error = errors.New("cart is empty")
)

// Cart represents user's shopping cart.
//
// Business rules:
//   - Items can be modified only while cart is active.
//   - Each variant can exist only once.
//   - Quantity must be > 0.
//   - Total is calculated from price snapshots.
type Cart struct {
	BaseAggregate

	id      int
	userID  int
	items   []CartItem
	status  CartStatus
	version int
}

// CartItem represents one product variant in cart.
type CartItem struct {
	variantID int
	quantity  int
	price     int64
}

// NewCart creates new active cart.
func NewCart(userID int) (*Cart, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	c := &Cart{
		userID: userID,
		status: CartStatusActive,
	}

	c.setInitialVersion(1)

	return c, nil
}

// ---- GETTERS ----

// ID returns cart id.
func (c *Cart) ID() int {
	return c.id
}

// UserID returns owner id.
func (c *Cart) UserID() int {
	return c.userID
}

// Status returns cart status.
func (c *Cart) Status() CartStatus {
	return c.status
}

// Version returns aggregate version.
func (c *Cart) Version() int {
	return c.version
}

// Items returns copy of cart items.
func (c *Cart) Items() []CartItem {
	result := make([]CartItem, len(c.items))
	copy(result, c.items)
	return result
}

// AddItem adds new variant to cart.
func (c *Cart) AddItem(variantID int, quantity int, price int64) error {
	if c.status != CartStatusActive {
		return ErrCartNotActive
	}

	if variantID <= 0 || price <= 0 || quantity <= 0 {
		return ErrInvalidItemQuality
	}

	for _, item := range c.items {
		if item.variantID == variantID {
			return ErrItemAlreadyExists
		}
	}

	c.items = append(c.items, CartItem{
		variantID: variantID,
		quantity:  quantity,
		price:     price,
	})

	c.incrementVersion()
	return nil
}

// RemoveItem removes variant from cart.
func (c *Cart) RemoveItem(variantID int) error {
	if c.status != CartStatusActive {
		return ErrCartNotActive
	}

	for i, item := range c.items {
		if item.variantID == variantID {
			c.items = append(c.items[:i], c.items[i+1:]...)

			c.incrementVersion()
			return nil
		}
	}

	return ErrItemNotFound
}

// ChangeQuantity updates item quantity.
func (c *Cart) ChangeQuantity(variantID int, quantity int) error {
	if c.status != CartStatusActive {
		return ErrCartNotActive
	}

	if quantity <= 0 {
		return ErrInvalidItemQuality
	}

	for i := range c.items {
		if c.items[i].variantID == variantID {
			c.items[i].quantity = quantity
			c.incrementVersion()
			return nil
		}
	}

	return ErrItemNotFound
}

// Total calculates total cart amount.
func (c *Cart) Total() int64 {
	var total int64
	for _, item := range c.items {
		total += int64(item.quantity) * item.price
	}
	return total
}

// Checkout closes cart.
//
// Fails if cart is empty or already checked out.
func (c *Cart) Checkout() error {
	if c.status != CartStatusActive {
		return ErrCartNotActive
	}

	if len(c.items) == 0 {
		return ErrCartEmpty
	}

	c.status = CartStatusCheckedOut
	c.incrementVersion()
	return nil
}
