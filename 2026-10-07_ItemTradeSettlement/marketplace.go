package marketplace

import (
	"errors"
	"fmt"
)

// Errors returned by the store and the trade service.
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

// Store is an in-memory repository of players. In production this would be
// backed by a real database; here it stands in for one so the exercise can
// run without any external dependencies.
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

// RemoveItem takes itemID out of playerID's inventory.
func (s *Store) RemoveItem(playerID, itemID string) error {
	p, err := s.GetPlayer(playerID)
	if err != nil {
		return err
	}
	for i, it := range p.Inventory {
		if it == itemID {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}

// AddItem places itemID into playerID's inventory.
func (s *Store) AddItem(playerID, itemID string) error {
	p, err := s.GetPlayer(playerID)
	if err != nil {
		return err
	}
	p.Inventory = append(p.Inventory, itemID)
	return nil
}

// AddGold credits playerID's balance.
func (s *Store) AddGold(playerID string, amount int) error {
	p, err := s.GetPlayer(playerID)
	if err != nil {
		return err
	}
	p.Gold += amount
	return nil
}

// RemoveGold debits playerID's balance, failing if the balance is too low.
func (s *Store) RemoveGold(playerID string, amount int) error {
	p, err := s.GetPlayer(playerID)
	if err != nil {
		return err
	}
	if p.Gold < amount {
		return ErrInsufficientGold
	}
	p.Gold -= amount
	return nil
}

// TradeService settles item-for-gold trades between two players.
type TradeService struct {
	store *Store
}

// NewTradeService returns a TradeService backed by store.
func NewTradeService(store *Store) *TradeService {
	return &TradeService{store: store}
}

// ExecuteTrade transfers itemID from sellerID to buyerID in exchange for
// price gold. Both players must exist, the seller must own the item, and
// the buyer must be able to afford it.
func (s *TradeService) ExecuteTrade(sellerID, buyerID, itemID string, price int) error {
	seller, err := s.store.GetPlayer(sellerID)
	if err != nil {
		return fmt.Errorf("seller: %w", err)
	}
	buyer, err := s.store.GetPlayer(buyerID)
	if err != nil {
		return fmt.Errorf("buyer: %w", err)
	}

	hasItem := false
	for _, it := range seller.Inventory {
		if it == itemID {
			hasItem = true
			break
		}
	}
	if !hasItem {
		return fmt.Errorf("seller %s: %w", sellerID, ErrItemNotFound)
	}
	if buyer.Gold < price {
		return fmt.Errorf("buyer %s: %w", buyerID, ErrInsufficientGold)
	}

	if err := s.store.RemoveItem(sellerID, itemID); err != nil {
		return fmt.Errorf("remove item from seller: %w", err)
	}
	if err := s.store.AddGold(sellerID, price); err != nil {
		return fmt.Errorf("credit seller: %w", err)
	}
	if err := s.store.RemoveGold(buyerID, price); err != nil {
		return fmt.Errorf("debit buyer: %w", err)
	}
	if err := s.store.AddItem(buyerID, itemID); err != nil {
		return fmt.Errorf("add item to buyer: %w", err)
	}

	return nil
}
