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

func (s *MemoryStore) RedeemAndCredit(code, userID string) (int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	card, ok := s.cards[code]
	if !ok {
		return 0, 0, ErrCardNotFound
	}
	if card.Redeemed {
		return 0, 0, ErrCardRedeemed
	}
	if s.CreditErr != nil {
		return 0, 0, s.CreditErr
	}

	card.Redeemed = true
	card.RedeemedBy = userID
	s.cards[code] = card
	s.balances[userID] += card.ValueCents

	return card.ValueCents, s.balances[userID], nil
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

	amountCents, newBalance, err := s.store.RedeemAndCredit(cardCode, userID)
	if err != nil {
		return Receipt{}, fmt.Errorf("redeem gift card: %w", err)
	}

	return Receipt{
		UserID:      userID,
		CardCode:    cardCode,
		AmountCents: amountCents,
		NewBalance:  newBalance,
	}, nil
}
