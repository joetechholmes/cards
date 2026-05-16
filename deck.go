package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Card struct {
	Suit  string
	Value string
}

func (c Card) String() string {
	return c.Value + " of " + c.Suit
}

type deck []Card

func newDeck() deck {
	cards := deck{}
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			cards = append(cards, Card{Suit: suit, Value: value})
		}
	}

	return cards
}

func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

// Returns independent copies to prevent shared backing-array aliasing.
func deal(d deck, handsize int) (deck, deck) {
	hand := make(deck, handsize)
	copy(hand, d[:handsize])
	rest := make(deck, len(d)-handsize)
	copy(rest, d[handsize:])
	return hand, rest
}

func (d deck) toString() string {
	strs := make([]string, len(d))
	for i, card := range d {
		strs[i] = card.Value + "|" + card.Suit
	}
	return strings.Join(strs, ",")
}

func (d deck) saveToFile(filename string) error {
	return os.WriteFile(filename, []byte(d.toString()), 0666)
}

func newDeckFromFile(filename string) deck {
	bs, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	parts := strings.Split(string(bs), ",")
	d := deck{}
	for _, p := range parts {
		halves := strings.Split(p, "|")
		if len(halves) == 2 {
			d = append(d, Card{Value: halves[0], Suit: halves[1]})
		}
	}
	return d
}

func (d deck) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPosition := r.Intn(len(d) - 1)
		d[i], d[newPosition] = d[newPosition], d[i]
	}
}
