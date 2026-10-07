package marketplace

import "testing"

// These helpers are expected to change as you refactor — update them to match your new API.
func newTestService(t *testing.T) (*TradeService, *Store) {
	t.Helper()
	store := NewStore()
	return NewTradeService(store), store
}

func seedPlayer(t *testing.T, store *Store, id string, gold int, items []string) {
	t.Helper()
	store.AddPlayer(&Player{ID: id, Gold: gold, Inventory: items})
}

func executeTrade(t *testing.T, svc *TradeService, sellerID, buyerID, itemID string, price int) error {
	t.Helper()
	return svc.ExecuteTrade(sellerID, buyerID, itemID, price)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestExecuteTrade(t *testing.T) {
	t.Run("seller not found", func(t *testing.T) {
		svc, store := newTestService(t)
		seedPlayer(t, store, "buyer1", 100, nil)

		t.Run("returns an error", func(t *testing.T) {
			err := executeTrade(t, svc, "ghost-seller", "buyer1", "sword", 50)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("buyer not found", func(t *testing.T) {
		svc, store := newTestService(t)
		seedPlayer(t, store, "seller1", 0, []string{"sword"})

		t.Run("returns an error", func(t *testing.T) {
			err := executeTrade(t, svc, "seller1", "ghost-buyer", "sword", 50)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	})

	t.Run("seller does not own the item", func(t *testing.T) {
		svc, store := newTestService(t)
		seedPlayer(t, store, "seller1", 0, []string{"shield"})
		seedPlayer(t, store, "buyer1", 100, nil)

		err := executeTrade(t, svc, "seller1", "buyer1", "sword", 50)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})

		t.Run("does not change the seller's inventory", func(t *testing.T) {
			p, _ := store.GetPlayer("seller1")
			if len(p.Inventory) != 1 || p.Inventory[0] != "shield" {
				t.Fatalf("seller inventory changed unexpectedly: %v", p.Inventory)
			}
		})
	})

	t.Run("buyer cannot afford the item", func(t *testing.T) {
		svc, store := newTestService(t)
		seedPlayer(t, store, "seller1", 0, []string{"sword"})
		seedPlayer(t, store, "buyer1", 10, nil)

		err := executeTrade(t, svc, "seller1", "buyer1", "sword", 50)

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})

		t.Run("does not remove the item from the seller", func(t *testing.T) {
			p, _ := store.GetPlayer("seller1")
			if len(p.Inventory) != 1 {
				t.Fatalf("seller lost the item despite a failed trade: %v", p.Inventory)
			}
		})
	})

	t.Run("successful trade", func(t *testing.T) {
		svc, store := newTestService(t)
		seedPlayer(t, store, "seller1", 0, []string{"sword"})
		seedPlayer(t, store, "buyer1", 100, nil)

		err := executeTrade(t, svc, "seller1", "buyer1", "sword", 50)

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})

		t.Run("removes the item from the seller", func(t *testing.T) {
			p, _ := store.GetPlayer("seller1")
			if len(p.Inventory) != 0 {
				t.Fatalf("expected seller to have no items, got %v", p.Inventory)
			}
		})

		t.Run("credits the seller with the price", func(t *testing.T) {
			p, _ := store.GetPlayer("seller1")
			if p.Gold != 50 {
				t.Fatalf("expected seller gold 50, got %d", p.Gold)
			}
		})

		t.Run("debits the buyer by the price", func(t *testing.T) {
			p, _ := store.GetPlayer("buyer1")
			if p.Gold != 50 {
				t.Fatalf("expected buyer gold 50, got %d", p.Gold)
			}
		})

		t.Run("adds the item to the buyer", func(t *testing.T) {
			p, _ := store.GetPlayer("buyer1")
			if len(p.Inventory) != 1 || p.Inventory[0] != "sword" {
				t.Fatalf("expected buyer to have the sword, got %v", p.Inventory)
			}
		})
	})
}
