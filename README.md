# cards

A terminal-based card game program written in Go. Currently supports Go Fish, with Solitaire and Poker planned.

## Games

### Go Fish
Play Go Fish against the computer. The goal is to collect the most "books" — sets of all four suits for a given rank (e.g., all four Aces).

**Rules:**
- Each player starts with 7 cards
- On your turn, ask for a rank you already hold (e.g. `Ace`, `King`, `Ten`)
- If the computer has that rank, it hands over all matching cards
- If not, you "Go Fish" and draw from the deck
- When you collect all four suits of a rank, it's removed as a completed book
- Game ends when all 13 books are claimed — most books wins

## Running

```
go run .
```

## Testing

```
go test ./...
```

## Project Structure

| File | Purpose |
|---|---|
| `main.go` | Entry point |
| `deck.go` | `Card` struct, `deck` type, shuffle, deal, save/load |
| `gofish.go` | Go Fish game loop and helper functions |
| `deck_test.go` | Tests for deck operations and game helpers |
