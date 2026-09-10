package cache

import (
	"fmt"
	"time"
)

const CartItemsCacheTTL = 10 * time.Minute

func CartItemsKey(cartID int64) string {
	return fmt.Sprintf("cart:items:%d", cartID)
}
