package marketplace

import (
	"errors"
	"fmt"
)

// Errors returned while looking up players or validating a trade.
var (
	ErrPlayerNotFound   = errors.New("player not found")
	ErrItemNotFound     = errors.New("item not found in inventory")
	ErrInsufficientGold = errors.New("insufficient gold")
)

// Player is an account in the marketplace economy.
type Player struct {
	ID        string
	Gold      int
	Inventory []string
}

// Store is an in-memory lookup of players, standing in for the real
// database this would hit in production so the exercise can run without
// any external dependencies.
type Store struct {
	players map[string]*Player
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{players: make(map[string]*Player)}
}

// AddPlayer registers a player in the store.
func (s *Store) AddPlayer(p *Player) {
	s.players[p.ID] = p
}

// GetPlayer looks up a player by ID.
func (s *Store) GetPlayer(id string) (*Player, error) {
	p, ok := s.players[id]
	if !ok {
		return nil, ErrPlayerNotFound
	}
	return p, nil
}

// ExecuteTrade transfers itemID from sellerID to buyerID in exchange for
// price gold. Both players must exist, the seller must own the item, and
// the buyer must be able to afford it.
func ExecuteTrade(store *Store, sellerID, buyerID, itemID string, price int) error {
	seller, err := store.GetPlayer(sellerID)
	if err != nil {
		return fmt.Errorf("seller: %w", err)
	}
	buyer, err := store.GetPlayer(buyerID)
	if err != nil {
		return fmt.Errorf("buyer: %w", err)
	}

	itemIndex := -1
	for i, it := range seller.Inventory {
		if it == itemID {
			itemIndex = i
			break
		}
	}
	if itemIndex == -1 {
		return fmt.Errorf("seller %s: %w", sellerID, ErrItemNotFound)
	}
	if buyer.Gold < price {
		return fmt.Errorf("buyer %s: %w", buyerID, ErrInsufficientGold)
	}

	// Step 1: take the item off the seller.
	seller.Inventory = append(seller.Inventory[:itemIndex], seller.Inventory[itemIndex+1:]...)
	// Step 2: pay the seller.
	seller.Gold += price
	// Step 3: charge the buyer.
	buyer.Gold -= price
	// Step 4: hand the item to the buyer.
	buyer.Inventory = append(buyer.Inventory, itemID)

	return nil
}
