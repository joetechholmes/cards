package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func (d deck) hasRank(rank string) bool {
	for _, card := range d {
		if card.Value == rank {
			return true
		}
	}
	return false
}

func removeRank(d *deck, rank string) deck {
	taken := deck{}
	remaining := deck{}
	for _, card := range *d {
		if card.Value == rank {
			taken = append(taken, card)
		} else {
			remaining = append(remaining, card)
		}
	}
	*d = remaining
	return taken
}

// Removes completed sets of 4 from hand and returns the number of books made.
func checkBooks(hand *deck) int {
	counts := make(map[string]int)
	for _, card := range *hand {
		counts[card.Value]++
	}
	isBook := make(map[string]bool)
	books := 0
	for value, count := range counts {
		if count == 4 {
			isBook[value] = true
			books++
		}
	}
	if books > 0 {
		remaining := deck{}
		for _, card := range *hand {
			if !isBook[card.Value] {
				remaining = append(remaining, card)
			}
		}
		*hand = remaining
	}
	return books
}

func drawCard(hand *deck, draw *deck) bool {
	if len(*draw) == 0 {
		return false
	}
	*hand = append(*hand, (*draw)[0])
	*draw = (*draw)[1:]
	return true
}

func playGoFish() {
	fmt.Println("=== Go Fish ===")
	fmt.Println("Match all four suits of a rank to complete a book. Most books wins!")

	draw := newDeck()
	draw.shuffle()

	var playerHand, computerHand deck
	playerHand, draw = deal(draw, 7)
	computerHand, draw = deal(draw, 7)

	playerBooks := 0
	computerBooks := 0

	playerBooks += checkBooks(&playerHand)
	computerBooks += checkBooks(&computerHand)

	scanner := bufio.NewScanner(os.Stdin)

	for playerBooks+computerBooks < 13 {
		// === Player's Turn ===
		if len(playerHand) == 0 {
			if !drawCard(&playerHand, &draw) {
				break
			}
			fmt.Println("\nYour hand was empty; you drew a card.")
		}

		fmt.Printf("\nYour hand (%d cards): ", len(playerHand))
		for _, c := range playerHand {
			fmt.Printf("[%s] ", c)
		}
		fmt.Printf("\nDraw pile: %d  |  Your books: %d  |  Computer books: %d\n", len(draw), playerBooks, computerBooks)

		var rank string
		for {
			fmt.Print("Ask for a rank (e.g. Ace, Ten, King): ")
			if !scanner.Scan() {
				return
			}
			rank = capitalize(strings.TrimSpace(scanner.Text()))
			if playerHand.hasRank(rank) {
				break
			}
			fmt.Printf("You don't have any %ss in your hand — you can only ask for a rank you hold. Try again.\n", rank)
		}

		if computerHand.hasRank(rank) {
			taken := removeRank(&computerHand, rank)
			playerHand = append(playerHand, taken...)
			fmt.Printf("Computer had %d %s(s)! You took them.\n", len(taken), rank)
		} else {
			fmt.Print("Go Fish! ")
			if drawCard(&playerHand, &draw) {
				drawn := playerHand[len(playerHand)-1]
				fmt.Printf("You drew: %s\n", drawn)
			} else {
				fmt.Println("The draw pile is empty.")
			}
		}

		if books := checkBooks(&playerHand); books > 0 {
			playerBooks += books
			fmt.Printf("You completed %d book(s)! Your total: %d\n", books, playerBooks)
		}

		if playerBooks+computerBooks == 13 {
			break
		}

		// === Computer's Turn ===
		if len(computerHand) == 0 {
			if !drawCard(&computerHand, &draw) {
				break
			}
		}

		randCard := computerHand[rand.Intn(len(computerHand))]
		askRank := randCard.Value
		fmt.Printf("\nComputer asks: Do you have any %ss?\n", askRank)

		if playerHand.hasRank(askRank) {
			taken := removeRank(&playerHand, askRank)
			computerHand = append(computerHand, taken...)
			fmt.Printf("You had %d %s(s)! Computer takes them.\n", len(taken), askRank)
		} else {
			fmt.Print("Go Fish! Computer draws a card")
			if drawCard(&computerHand, &draw) {
				fmt.Println(".")
			} else {
				fmt.Println(" — but the draw pile is empty!")
			}
		}

		if books := checkBooks(&computerHand); books > 0 {
			computerBooks += books
			fmt.Printf("Computer completed %d book(s)! Computer total: %d\n", books, computerBooks)
		}
	}

	fmt.Printf("\n=== Game Over ===\n")
	fmt.Printf("Your books: %d  |  Computer books: %d\n", playerBooks, computerBooks)
	switch {
	case playerBooks > computerBooks:
		fmt.Println("You win!")
	case computerBooks > playerBooks:
		fmt.Println("Computer wins!")
	default:
		fmt.Println("It's a tie!")
	}
}
