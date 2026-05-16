package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(d))
	}

	if d[0] != (Card{Value: "Ace", Suit: "Spades"}) {
		t.Errorf("Expected Ace of Spades, but got %v", d[0])
	}

	if d[len(d)-1] != (Card{Value: "King", Suit: "Clubs"}) {
		t.Errorf("Expected King of Clubs, but got %v", d[len(d)-1])
	}
}

func TestNewDeckContainsAllUniqueCards(t *testing.T) {
	d := newDeck()
	seen := make(map[string]bool)
	for _, card := range d {
		key := card.Value + "|" + card.Suit
		if seen[key] {
			t.Errorf("Duplicate card in deck: %v", card)
		}
		seen[key] = true
	}
	if len(seen) != 52 {
		t.Errorf("Expected 52 unique cards, got %v", len(seen))
	}
}

func TestCardString(t *testing.T) {
	c := Card{Value: "Ace", Suit: "Spades"}
	if c.String() != "Ace of Spades" {
		t.Errorf("Expected 'Ace of Spades', got %v", c.String())
	}
}

func TestDeal(t *testing.T) {
	d := newDeck()
	hand, rest := deal(d, 7)

	if len(hand) != 7 {
		t.Errorf("Expected hand size 7, got %v", len(hand))
	}
	if len(rest) != 45 {
		t.Errorf("Expected rest size 45, got %v", len(rest))
	}
	if hand[0] != (Card{Value: "Ace", Suit: "Spades"}) {
		t.Errorf("Expected first hand card to be Ace of Spades, got %v", hand[0])
	}
}

func TestDealIndependence(t *testing.T) {
	d := newDeck()
	hand, rest := deal(d, 7)

	hand = append(hand, Card{Value: "Joker", Suit: "None"})

	if len(rest) != 45 {
		t.Errorf("Appending to hand affected rest: expected rest size 45, got %v", len(rest))
	}
}

func TestShufflePreservesLength(t *testing.T) {
	d := newDeck()
	d.shuffle()
	if len(d) != 52 {
		t.Errorf("Expected 52 cards after shuffle, got %v", len(d))
	}
}

func TestShuffleChangesOrder(t *testing.T) {
	d := newDeck()
	original := make(deck, len(d))
	copy(original, d)
	d.shuffle()

	same := true
	for i := range d {
		if d[i] != original[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("Deck order unchanged after shuffle — shuffle may not be working")
	}
}

func TestSaveToDeckAndNewDeckFromFile(t *testing.T) {
	os.Remove("_decktesting")

	deck := newDeck()

	deck.saveToFile("_decktesting")

	loadedDeck := newDeckFromFile("_decktesting")

	if len(loadedDeck) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(loadedDeck))
	}

	os.Remove("_decktesting")
}

func TestSaveAndLoadPreservesCards(t *testing.T) {
	os.Remove("_decktesting2")

	original := newDeck()
	original.saveToFile("_decktesting2")
	loaded := newDeckFromFile("_decktesting2")

	for i, card := range original {
		if loaded[i] != card {
			t.Errorf("Card mismatch at index %d: expected %v, got %v", i, card, loaded[i])
		}
	}

	os.Remove("_decktesting2")
}

func TestHasRank(t *testing.T) {
	hand := deck{
		{Value: "Ace", Suit: "Spades"},
		{Value: "King", Suit: "Hearts"},
	}

	if !hand.hasRank("Ace") {
		t.Error("Expected hand to contain Ace")
	}
	if !hand.hasRank("King") {
		t.Error("Expected hand to contain King")
	}
	if hand.hasRank("Two") {
		t.Error("Expected hand not to contain Two")
	}
}

func TestRemoveRank(t *testing.T) {
	hand := deck{
		{Value: "Ace", Suit: "Spades"},
		{Value: "Ace", Suit: "Hearts"},
		{Value: "King", Suit: "Clubs"},
	}

	taken := removeRank(&hand, "Ace")

	if len(taken) != 2 {
		t.Errorf("Expected 2 cards taken, got %v", len(taken))
	}
	if len(hand) != 1 {
		t.Errorf("Expected 1 card remaining, got %v", len(hand))
	}
	if hand[0].Value != "King" {
		t.Errorf("Expected remaining card to be King, got %v", hand[0].Value)
	}

	// Removing a rank not present should leave hand unchanged
	taken = removeRank(&hand, "Two")
	if len(taken) != 0 {
		t.Errorf("Expected 0 cards taken for missing rank, got %v", len(taken))
	}
	if len(hand) != 1 {
		t.Errorf("Expected hand unchanged after removing missing rank, got %v", len(hand))
	}
}

func TestCheckBooks(t *testing.T) {
	hand := deck{
		{Value: "Ace", Suit: "Spades"},
		{Value: "Ace", Suit: "Hearts"},
		{Value: "Ace", Suit: "Diamonds"},
		{Value: "Ace", Suit: "Clubs"},
		{Value: "King", Suit: "Spades"},
	}

	books := checkBooks(&hand)

	if books != 1 {
		t.Errorf("Expected 1 book, got %v", books)
	}
	if len(hand) != 1 {
		t.Errorf("Expected 1 card remaining after book removed, got %v", len(hand))
	}
	if hand[0].Value != "King" {
		t.Errorf("Expected remaining card to be King, got %v", hand[0].Value)
	}
}

func TestCheckBooksNoBook(t *testing.T) {
	hand := deck{
		{Value: "Ace", Suit: "Spades"},
		{Value: "Ace", Suit: "Hearts"},
		{Value: "King", Suit: "Clubs"},
	}

	books := checkBooks(&hand)

	if books != 0 {
		t.Errorf("Expected 0 books, got %v", books)
	}
	if len(hand) != 3 {
		t.Errorf("Expected hand unchanged at 3 cards, got %v", len(hand))
	}
}

func TestDrawCard(t *testing.T) {
	draw := deck{{Value: "Ace", Suit: "Spades"}, {Value: "King", Suit: "Hearts"}}
	hand := deck{}

	ok := drawCard(&hand, &draw)

	if !ok {
		t.Error("Expected drawCard to return true")
	}
	if len(hand) != 1 {
		t.Errorf("Expected hand size 1, got %v", len(hand))
	}
	if hand[0] != (Card{Value: "Ace", Suit: "Spades"}) {
		t.Errorf("Expected drawn card to be Ace of Spades, got %v", hand[0])
	}
	if len(draw) != 1 {
		t.Errorf("Expected draw pile size 1, got %v", len(draw))
	}
}

func TestDrawCardEmptyPile(t *testing.T) {
	draw := deck{}
	hand := deck{}

	ok := drawCard(&hand, &draw)

	if ok {
		t.Error("Expected drawCard to return false on empty draw pile")
	}
	if len(hand) != 0 {
		t.Errorf("Expected hand to remain empty, got %v", len(hand))
	}
}
