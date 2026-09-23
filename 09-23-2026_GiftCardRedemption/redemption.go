package redemption

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrCardNotFound = errors.New("gift card not found")
	ErrCardRedeemed = errors.New("gift card already redeemed")
	ErrUnavailable  = errors.New("store unavailable")
)

type GiftCard struct {
	Code       string
	ValueCents int
	Redeemed   bool
	RedeemedBy string
}

type Receipt struct {
	UserID      string
	CardCode    string
	AmountCents int
	NewBalance  int
}

type MemoryStore struct {
	mu        sync.Mutex
	cards     map[string]GiftCard
	balances  map[string]int
	CreditErr error
}

func NewMemoryStore(cards ...GiftCard) *MemoryStore {
	store := &MemoryStore{
		cards:    make(map[string]GiftCard, len(cards)),
		balances: make(map[string]int),
	}
	for _, card := range cards {
		store.cards[card.Code] = card
	}
	return store
}

func (s *MemoryStore) FindGiftCard(code string) (GiftCard, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	card, ok := s.cards[code]
	if !ok {
		return GiftCard{}, ErrCardNotFound
	}
	return card, nil
}

func (s *MemoryStore) MarkRedeemed(code, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	card, ok := s.cards[code]
	if !ok {
		return ErrCardNotFound
	}
	if card.Redeemed {
		return ErrCardRedeemed
	}

	card.Redeemed = true
	card.RedeemedBy = userID
	s.cards[code] = card
	return nil
}

func (s *MemoryStore) CreditBalance(userID string, amountCents int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.CreditErr != nil {
		return 0, s.CreditErr
	}

	s.balances[userID] += amountCents
	return s.balances[userID], nil
}

func (s *MemoryStore) GiftCard(code string) GiftCard {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.cards[code]
}

func (s *MemoryStore) Balance(userID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.balances[userID]
}

type RedemptionService struct {
	store *MemoryStore
}

func NewRedemptionService(store *MemoryStore) *RedemptionService {
	return &RedemptionService{store: store}
}

func (s *RedemptionService) Redeem(userID, cardCode string) (Receipt, error) {
	userID = strings.TrimSpace(userID)
	cardCode = strings.ToUpper(strings.TrimSpace(cardCode))

	if userID == "" {
		return Receipt{}, ErrInvalidUser
	}

	card, err := s.store.FindGiftCard(cardCode)
	if err != nil {
		return Receipt{}, fmt.Errorf("find gift card: %w", err)
	}
	if card.Redeemed {
		return Receipt{}, ErrCardRedeemed
	}

	if err := s.store.MarkRedeemed(card.Code, userID); err != nil {
		return Receipt{}, fmt.Errorf("mark gift card redeemed: %w", err)
	}

	newBalance, err := s.store.CreditBalance(userID, card.ValueCents)
	if err != nil {
		return Receipt{}, fmt.Errorf("credit balance: %w", err)
	}

	return Receipt{
		UserID:      userID,
		CardCode:    card.Code,
		AmountCents: card.ValueCents,
		NewBalance:  newBalance,
	}, nil
}
